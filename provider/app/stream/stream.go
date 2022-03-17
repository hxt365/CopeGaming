package stream

import (
	"container/ring"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"provider/app/webrtc"
	"provider/constants"

	"github.com/pion/rtp"
)

type StreamRelayer struct {
	logID         string
	videoStream   chan *rtp.Packet
	audioStream   chan *rtp.Packet
	eventStream   chan *webrtc.Packet
	videoListener *net.UDPConn
	audioListener *net.UDPConn
	wineConn      *net.TCPConn
	syncListener  *net.TCPListener
}

func NewStreamRelayer(logID string, videoStream, audioStream chan *rtp.Packet, eventStream chan *webrtc.Packet, videoListener, audioListener *net.UDPConn, syncListener *net.TCPListener) *StreamRelayer {
	s := &StreamRelayer{
		logID:         logID,
		videoStream:   videoStream,
		audioStream:   audioStream,
		eventStream:   eventStream,
		videoListener: videoListener,
		audioListener: audioListener,
		syncListener:  syncListener,
	}

	return s
}

func (s *StreamRelayer) Start() error {
	log.Printf("[%s] Start relaying streams..\n", s.logID)

	wineConnected := make(chan struct{})

	go func() {
		for {
			log.Printf("[%s] Waiting for syncinput to connect\n", s.logID)

			conn, err := s.syncListener.AcceptTCP()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("[%s] Couldn't accept syncinput connection: %s\n", s.logID, err)
				continue
			}

			err = conn.SetKeepAlive(true)
			if err != nil {
				log.Printf("[%s] Couldn't set keepAlive: %s\n", s.logID, err)
				continue
			}
			err = conn.SetKeepAlivePeriod(10 * time.Second)
			if err != nil {
				log.Printf("[%s] Couldn't set keepAlive period: %s\n", s.logID, err)
				continue
			}

			if s.wineConn != nil {
				_ = s.wineConn.Close()
			}
			s.wineConn = conn
			wineConnected <- struct{}{}

			log.Printf("[%s] Successfully set up syncinput connection\n", s.logID)
		}
	}()

	go func() {
		<-wineConnected

		go s.healthCheckVM()
		go s.handleAppEvents()
		go s.relayStream(s.videoListener, s.videoStream)
		go s.relayStream(s.audioListener, s.audioStream)
	}()

	return nil
}

func (s *StreamRelayer) Close() {
	if s.wineConn != nil {
		_ = s.wineConn.Close()
	}
}

func (s *StreamRelayer) relayStream(listener *net.UDPConn, output chan<- *rtp.Packet) {
	r := ring.New(120)

	n := r.Len()
	for i := 0; i < n; i++ {
		r.Value = make([]byte, 1500)
		r = r.Next()
	}

	for {
		inboundRTPPacket := r.Value.([]byte)
		r = r.Next()

		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("[%s] Error during read RTP packet: %s\n", s.logID, err)
			continue
		}

		var packet rtp.Packet
		if err := packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
			log.Printf("[%s] Error during unmarshalling RTP packet: %s\n", s.logID, err)
			continue
		}

		output <- &packet
	}
}

func (s *StreamRelayer) healthCheckVM() {
	for {
		if s.wineConn != nil {
			_, err := s.wineConn.Write([]byte{0})
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("[%s] Error while health checking: %s\n", s.logID, err)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func (s *StreamRelayer) handleAppEvents() {
	for packet := range s.eventStream {
		switch packet.Type {
		case constants.KeyUp:
			s.simulateKey(packet.Data, 0)
		case constants.KeyDown:
			s.simulateKey(packet.Data, 1)
		case constants.MouseMove:
			s.simulateMouseEvent(packet.Data, 0)
		case constants.MouseDown:
			s.simulateMouseEvent(packet.Data, 1)
		case constants.MouseUp:
			s.simulateMouseEvent(packet.Data, 2)
		}
	}
}

type keydownPayload struct {
	KeyCode int `json:"keycode"`
}

func (s *StreamRelayer) simulateKey(jsonPayload string, keyState byte) {
	if s.wineConn == nil {
		return
	}

	p := &keydownPayload{}
	err := json.Unmarshal([]byte(jsonPayload), &p)
	if err != nil {
		log.Printf("[%s] Couldn't parse keydown payload: %s\n", s.logID, err)
		return
	}

	vmKeyMsg := fmt.Sprintf("K%d,%b|", p.KeyCode, keyState)
	_, err = s.wineConn.Write([]byte(vmKeyMsg))
	if err != nil {
		log.Printf("[%s] Couldn't send key event to wine: %s\n", s.logID, err)
		return
	}
}

type mousePayload struct {
	IsLeft byte    `json:"isLeft"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func (s *StreamRelayer) simulateMouseEvent(jsonPayload string, mouseState int) {
	if s.wineConn == nil {
		return
	}

	p := &mousePayload{}
	err := json.Unmarshal([]byte(jsonPayload), &p)
	if err != nil {
		log.Printf("[%s] Couldn't parse mouse payload: %s\n", s.logID, err)
		return
	}

	// Original formula: p.X = p.X * s.screenWidth / p.Width
	// However, syncinput will handle the multiplication with screenWidth to remove the need of configuring screenWidth
	// Same for p.Y
	p.X = p.X / p.Width
	p.Y = p.Y / p.Height

	// Mouse is in format of comma separated "12.4,52.3"
	vmMouseMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, mouseState, p.X, p.Y, p.Width, p.Height)
	_, err = s.wineConn.Write([]byte(vmMouseMsg))
	if err != nil {
		log.Printf("[%s] Couldn't send mouse event to wine: %s\n", s.logID, err)
		return
	}
}
