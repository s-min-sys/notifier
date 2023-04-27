package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/s-min-sys/notifier/internal/config"
	"github.com/s-min-sys/notifier/internal/inters"
	"github.com/s-min-sys/notifier/internal/telebot"
	"github.com/s-min-sys/notifier/pkg"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
)

func NewServer(logger l.Wrapper) *Server {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	s := &Server{
		logger: logger.WithFields(l.StringField(l.ClsKey, "server")),
	}

	s.init()

	return s
}

type Server struct {
	wg     sync.WaitGroup
	logger l.Wrapper

	teleNotifier inters.Sender
}

func (s *Server) init() {
	cfg := config.GetConfig()

	if cfg.TeleConfig != nil {
		s.teleNotifier = telebot.NewSender(cfg.TeleConfig.Token, cfg.TeleConfig.APIEndPoint, cfg.RedisCli, s.logger)
	}

	s.wg.Add(1)

	go s.httpServerRoutine()
}

func (s *Server) httpServerRoutine() {
	defer s.wg.Done()

	mux := http.NewServeMux()

	mux.HandleFunc("/send-text-message", func(writer http.ResponseWriter, request *http.Request) {
		d, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		var req pkg.TextMessage

		_ = json.Unmarshal(d, &req)
		if req.Text == "" {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = s.SendTextMessage(req)

		if err == nil {
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
		}
	})

	// 创建服务器
	server := &http.Server{
		Addr:         config.GetConfig().Listen,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	log.Fatal(server.ListenAndServe())
}

func (s *Server) SendTextMessage(message pkg.TextMessage) (err error) {
	var senderIDs []pkg.SenderID

	if message.SenderID == pkg.SenderIDAll {
		senderIDs = append(senderIDs, pkg.SenderIDTelegram)
		senderIDs = append(senderIDs, pkg.SenderIDWeChat)
	} else {
		senderIDs = append(senderIDs, message.SenderID)
	}

	for _, senderID := range senderIDs {
		curMessage := message
		curMessage.SenderID = senderID

		switch senderID {
		case pkg.SenderIDTelegram:
			if s.teleNotifier == nil {
				err = commerr.ErrUnavailable
			} else {
				err = s.teleNotifier.SendTextMessage(curMessage)
			}
		case pkg.SenderIDWeChat:
			err = commerr.ErrUnimplemented
		default:
			err = commerr.ErrInvalidArgument
		}

		if err != nil {
			s.logger.WithFields(l.ErrorField(err)).Error("send to %v failed", senderID)
		}
	}

	return
}

func (s *Server) Wait() {
	s.wg.Wait()
}
