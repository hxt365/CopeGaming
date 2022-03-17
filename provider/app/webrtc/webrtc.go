package webrtc

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"

	"provider/pkg/socket"
	"provider/settings"
	"provider/utils"

	"github.com/pion/interceptor"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type WebRTC struct {
	logID        string
	conn         *webrtc.PeerConnection
	imageChannel chan *rtp.Packet
	audioChannel chan *rtp.Packet
	eventChannel chan *Packet
	inputTrack   *webrtc.DataChannel
	healthTrack  *webrtc.DataChannel
	// Connection close signal channel
	closed       chan struct{}
	// Ensure that exit callback function will be called only once
	exitOnce     sync.Once
}

type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type OnIceCallback func(candidate string)
type OnExitCallback func()

var (
	webrtcSettings webrtc.SettingEngine
	settingOnce    sync.Once
)

const MaxMissedHealthCheck int = 5

func NewWebRTC(logID string, videoStream, audioStream chan *rtp.Packet, inputStream chan *Packet) (*WebRTC, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	i := &interceptor.Registry{}
	if !settings.DisableDefaultInterceptors {
		if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
			return nil, err
		}
	}

	settingOnce.Do(func() {
		settingEngine := webrtc.SettingEngine{}

		if settings.PortRange.Min > 0 && settings.PortRange.Max > 0 {
			if err := settingEngine.SetEphemeralUDPPortRange(settings.PortRange.Min, settings.PortRange.Max); err != nil {
				panic(err)
			}
		} else if settings.SinglePort > 0 {
			l, err := socket.NewSocketPortRoll("udp", settings.SinglePort)
			if err != nil {
				panic(err)
			}
			udpListener := l.(*net.UDPConn)
			log.Printf("[%s] Listening for WebRTC traffic at %s\n", logID, udpListener.LocalAddr())
			settingEngine.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpListener))
		}

		if settings.IceIpMap != "" {
			settingEngine.SetNAT1To1IPs([]string{settings.IceIpMap}, webrtc.ICECandidateTypeHost)
		}

		webrtcSettings = settingEngine
	})

	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(m),
		webrtc.WithInterceptorRegistry(i),
		webrtc.WithSettingEngine(webrtcSettings),
	)

	conn, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		}},
	)
	if err != nil {
		return nil, err
	}

	return &WebRTC{
		logID:        logID,
		conn:         conn,
		imageChannel: videoStream,
		audioChannel: audioStream,
		eventChannel: inputStream,
		closed:       make(chan struct{}),
	}, nil
}

func (w *WebRTC) StartClient(vCodec string, iceCb OnIceCallback, exitCb OnExitCallback) (string, error) {
	log.Printf("[%s] Start WebRTC..\n", w.logID)

	videoTrack, err := w.addVideoTrack(vCodec)
	if err != nil {
		return "", err
	}

	audioTrack, err := w.addAudioTrack()
	if err != nil {
		return "", err
	}

	err = w.addInputTrack(true)
	if err != nil {
		return "", err
	}

	err = w.addHealthCheck(exitCb)
	if err != nil {
		return "", err
	}

	w.conn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionStateConnected {
			log.Printf("[%s] ICE Connected succeeded\n", w.logID)
			w.startStreamingVideo(videoTrack)
			w.startStreamingAudio(audioTrack)
		}

		if state == webrtc.ICEConnectionStateFailed || state == webrtc.ICEConnectionStateClosed || state == webrtc.ICEConnectionStateDisconnected {
			log.Printf("[%s] ICE Connected failed: %s\n", w.logID, state)
			w.exitOnce.Do(exitCb)
		}
	})

	w.conn.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		if iceCandidate != nil {
			candidate, err := utils.EncodeBase64(iceCandidate.ToJSON())
			if err != nil {
				log.Printf("[%s] Encode IceCandidate failed: %s\n", w.logID, err)
				return
			}
			iceCb(candidate)
		}
	})

	// Create offer
	offer, err := w.conn.CreateOffer(nil)
	if err != nil {
		return "", err
	}

	err = w.conn.SetLocalDescription(offer)
	if err != nil {
		return "", err
	}

	encodedOffer, err := utils.EncodeBase64(offer)
	if err != nil {
		return "", err
	}

	return encodedOffer, nil
}

func (w *WebRTC) addVideoTrack(vCodec string) (*webrtc.TrackLocalStaticRTP, error) {
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: verbalCodecToMime(vCodec),
	}, "video", "pion")
	if err != nil {
		return nil, err
	}

	videoSender, err := w.conn.AddTrack(videoTrack)
	if err != nil {
		return nil, err
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := videoSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	return videoTrack, nil
}

func (w *WebRTC) addAudioTrack() (*webrtc.TrackLocalStaticRTP, error) {
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeOpus,
	}, "audio", "pion")
	if err != nil {
		return nil, err
	}

	audioSender, err := w.conn.AddTrack(audioTrack)
	if err != nil {
		return nil, err
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := audioSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	return audioTrack, nil
}

func (w *WebRTC) addInputTrack(unreliable bool) error {
	var options *webrtc.DataChannelInit
	if unreliable {
		f := false
		zero := uint16(0)
		options = &webrtc.DataChannelInit{
			Ordered:        &f,
			MaxRetransmits: &zero,
		}
	}

	inputTrack, err := w.conn.CreateDataChannel("app-input", options)
	if err != nil {
		return err
	}
	w.inputTrack = inputTrack

	inputTrack.OnMessage(func(rawMsg webrtc.DataChannelMessage) {
		defer func() {
			if r := recover(); r != nil {
				// Maybe sent to closed channel
			}
		}()

		var msg Packet
		if err := json.Unmarshal(rawMsg.Data, &msg); err != nil {
			log.Printf("[%s] Couldn't parse webrtc data message: %s\n", w.logID, err)
			return
		}

		w.eventChannel <- &msg
	})
	return nil
}

func (w *WebRTC) addHealthCheck(exitCb OnExitCallback) error {
	healthTrack, err := w.conn.CreateDataChannel("health-check", nil)
	if err != nil {
		return err
	}
	w.healthTrack = healthTrack

	missedHealthCheckCounts := 0
	lock := sync.Mutex{}

	go func() {
		for {
			select {
			case <-w.closed:
				return
			default:
				lock.Lock()
				missedHealthCheckCounts += 1
				if missedHealthCheckCounts == MaxMissedHealthCheck {
					log.Printf("[%s] Health-check failed", w.logID)
					w.exitOnce.Do(exitCb)
					return
				}
				lock.Unlock()
				time.Sleep(2 * time.Second)
			}
		}
	}()

	healthTrack.OnMessage(func(_ webrtc.DataChannelMessage) {
		lock.Lock()
		missedHealthCheckCounts = 0
		lock.Unlock()
	})

	return nil
}

func (w *WebRTC) SetRemoteSDP(remoteSDP string) error {
	var answer webrtc.SessionDescription

	err := utils.DecodeBase64(remoteSDP, &answer)
	if err != nil {
		log.Printf("[%s] Decode remote sdp from peer failed: %s\n", w.logID, err)
		return err
	}

	err = w.conn.SetRemoteDescription(answer)
	if err != nil {
		log.Printf("[%s] Set remote description from peer failed: %s\n", w.logID, err)
		return err
	}

	return nil
}

func (w *WebRTC) AddCandidate(candidate string) error {
	var iceCandidate webrtc.ICECandidateInit

	err := utils.DecodeBase64(candidate, &iceCandidate)
	if err != nil {
		log.Printf("[%s] Decode Ice candidate from peer failed: %s\n", w.logID, err)
		return err
	}

	err = w.conn.AddICECandidate(iceCandidate)
	if err != nil {
		log.Printf("[%s] Add Ice candidate from peer failed: %s\n", w.logID, err)
		return err
	}

	return nil
}

func (w *WebRTC) StopClient() {
	w.inputTrack.Close()
	w.healthTrack.Close()
	w.conn.Close()
	w.closed <- struct{}{}
}

func (w *WebRTC) startStreamingVideo(videoTrack *webrtc.TrackLocalStaticRTP) {
	go func() {
		for packet := range w.imageChannel {
			if err := videoTrack.WriteRTP(packet); err != nil {
				log.Printf("[%s] Error when writing RTP to video track: %s\n", w.logID, err)
			}
		}
	}()
}

func (w *WebRTC) startStreamingAudio(audioTrack *webrtc.TrackLocalStaticRTP) {
	go func() {
		for packet := range w.audioChannel {
			if err := audioTrack.WriteRTP(packet); err != nil {
				log.Printf("[%s] Error when writing RTP to opus track: %s\n", w.logID, err)
			}
		}
	}()
}

func verbalCodecToMime(vCodec string) string {
	switch vCodec {
	case "h264":
		return webrtc.MimeTypeH264
	case "vpx":
		return webrtc.MimeTypeVP8
	default:
		return webrtc.MimeTypeVP8
	}
}
