package session

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"provider/app/stream"
	"provider/app/vm"
	"provider/app/webrtc"
	"provider/app/ws"
	"provider/constants"
	"provider/pkg/socket"
	"provider/settings"
	"provider/utils"

	"github.com/pion/rtp"
)

type Session struct {
	playerID  string
	timeStart time.Time
	hub       *Hub
	// Inbound message buffer
	inpBuf chan *ws.Message
	// Outbound message buffer
	outBuf chan interface{}
	// WS connection to coordinator service
	wsConn *ws.Connection
}

func NewSession(playerID string, wsConn *ws.Connection, hub *Hub) *Session {
	s := Session{
		playerID:  playerID,
		hub:       hub,
		timeStart: time.Now(),
		inpBuf:    make(chan *ws.Message),
		outBuf:    make(chan interface{}),
		wsConn:    wsConn,
	}

	go s.readMsg()
	go s.writeMsg()

	return &s
}

func (s *Session) close() {
	defer func() {
		if r := recover(); r != nil {
			// maybe close already closed channel
		}
	}()

	close(s.inpBuf)
	close(s.outBuf)

	s.hub.RemoveSession(s.playerID)
}

func (s *Session) ReceiveMsg(msg *ws.Message) {
	s.inpBuf <- msg
}

func (s *Session) writeMsg() {
	for msg := range s.outBuf {
		if err := s.wsConn.Send(msg); err != nil {
			log.Printf("[%s] Failed to write message %s: %s", s.playerID, msg, err)
		}
	}
}

func (s *Session) sendIceCandidate(candidate string) {
	s.outBuf <- ws.Message{
		ReceiverID: s.playerID,
		Type:       constants.IceCandidateMessage,
		Data:       candidate,
	}
}

func (s *Session) sendOffer(offer string) {
	s.outBuf <- ws.Message{
		ReceiverID: s.playerID,
		Type:       constants.SDPMessage,
		Data:       offer,
	}
}

type Configure struct {
	Device string `json:"device"`
	AppID  string `json:"appID"`
}

func (s *Session) start(conf *Configure) (*webrtc.WebRTC, error) {
	// Create relaying streams
	videoStream := make(chan *rtp.Packet, 100)
	audioStream := make(chan *rtp.Packet, 100)
	inputStream := make(chan *webrtc.Packet, 100)

	videoListener, err := socket.NewRandomUDPListener()
	if err != nil {
		log.Printf("[%s] Couldn't create a UDP listener for video: %s\n", s.playerID, err)
		return nil, err
	}
	videoRelayPort, err := socket.ExtractPort(videoListener.LocalAddr().String())
	if err != nil {
		log.Printf("[%s] Couldn't extract UDP port for video: %s\n", s.playerID, err)
		return nil, err
	}
	audioListener, err := socket.NewRandomUDPListener()
	if err != nil {
		log.Printf("[%s] Couldn't create a UDP listener for audio: %s\n", s.playerID, err)
		return nil, err
	}
	audioRelayPort, err := socket.ExtractPort(audioListener.LocalAddr().String())
	if err != nil {
		log.Printf("[%s] Couldn't extract UDP port for audio: %s\n", s.playerID, err)
		return nil, err
	}
	syncListener, err := socket.NewRandomTCPListener()
	if err != nil {
		log.Printf("[%s] Couldn't create a TCP listener for wine: %s\n", s.playerID, err)
		return nil, err
	}
	syncPort, err := socket.ExtractPort(syncListener.Addr().String())
	if err != nil {
		log.Printf("[%s] Couldn't extract TCP port for wine: %s\n", s.playerID, err)
		return nil, err
	}

	log.Printf("[%s] Wait for video at port %d\n", s.playerID, videoRelayPort)
	log.Printf("[%s] Wait for audio at port %d\n", s.playerID, audioRelayPort)
	log.Printf("[%s] Wait for syncinput at port %d\n", s.playerID, syncPort)

	relayer := stream.NewStreamRelayer(s.playerID,
		videoStream, audioStream, inputStream,
		videoListener, audioListener, syncListener)
	if err := relayer.Start(); err != nil {
		fmt.Printf("[%s] Couldn't start relaying streams: %s\n", s.playerID, err)
		return nil, err
	}

	// Start VM
	appName := fmt.Sprintf("%s_%s", conf.AppID, conf.Device)
	appId := fmt.Sprintf("%s_%s", s.playerID, utils.RandString(6))
	if err := vm.StartVM(appId, appName, videoRelayPort, audioRelayPort, syncPort); err != nil {
		log.Printf("[%s] Error when start VM: %s\n", s.playerID, err)
		return nil, err
	}

	// Start WebRTC
	webrtcConn, err := webrtc.NewWebRTC(s.playerID, videoStream, audioStream, inputStream)
	if err != nil {
		return nil, err
	}

	onExitCb := func() {
		log.Printf("[%s] Releasing allocated resources", s.playerID)
		if err := vm.StopVM(appId, appName); err != nil {
			log.Printf("[%s] Error when stopping VM: %s\n", s.playerID, err)
		}

		// Must close webrtc connection first to ensure no writing to closed inputStream
		webrtcConn.StopClient()

		// Must close listeners before streams to ensure no writing to closed channels
		audioListener.Close()
		videoListener.Close()
		syncListener.Close()

		close(videoStream)
		close(audioStream)
		close(inputStream)

		relayer.Close()
		s.close()
	}
	offer, err := webrtcConn.StartClient(settings.VideoCodec, s.sendIceCandidate, onExitCb)
	if err != nil {
		fmt.Printf("[%s] Couldn't start webrtc client: %s\n", s.playerID, err)
		return nil, err
	}

	s.sendOffer(offer)

	return webrtcConn, nil
}

func (s *Session) readMsg() {
	var (
		webrtcConn *webrtc.WebRTC
		err        error
	)

	for msg := range s.inpBuf {
		switch msg.Type {
		case constants.StartMessage:
			var conf Configure
			if err := json.Unmarshal([]byte(msg.Data), &conf); err != nil {
				log.Printf("[%s] Error when parse Start message: %s\n", s.playerID, err)
				continue
			}
			webrtcConn, err = s.start(&conf)
			if err != nil {
				log.Printf("[%s] Error when starting new session: %s\n", s.playerID, err)
				webrtcConn = nil
			}
		case constants.SDPMessage:
			if webrtcConn == nil {
				continue
			}
			err := webrtcConn.SetRemoteSDP(msg.Data)
			if err != nil {
				log.Printf("[%s] Couldn't set remote SDP %s\n", s.playerID, err)
				webrtcConn = nil
			}
		case constants.IceCandidateMessage:
			if webrtcConn == nil {
				continue
			}
			err := webrtcConn.AddCandidate(msg.Data)
			if err != nil {
				log.Printf("[%s] Couldn't set ICE candidate %s\n", s.playerID, err)
			}
		}
	}
}
