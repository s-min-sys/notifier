package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	sharepkg "github.com/s-min-sys/notifier-share/pkg"
	"github.com/s-min-sys/notifier/internal/config"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/ptl"
)

func NewServer(cfg *config.Config, logger l.Wrapper) *Server {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	if cfg == nil {
		logger.Fatal("no config")
	}

	s := &Server{
		cfg:    cfg,
		logger: logger.WithFields(l.StringField(l.ClsKey, "server")),
	}

	s.init()

	return s
}

type Server struct {
	wg     sync.WaitGroup
	cfg    *config.Config
	logger l.Wrapper

	senders sync.Map
}

func (s *Server) init() {
	s.wg.Add(1)

	go s.httpServerRoutine()
}

func (s *Server) httpServerRoutine() {
	defer s.wg.Done()

	sharepkg.RunCommandServer[sharepkg.TextMessage](s.cfg.Listens, "", s, s.logger)
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) SendTextMessage(req *sharepkg.TextMessage) (code ptl.Code, msg string) {
	httpClient, senderURL, msg := s.getOrCreateClientForSenderID(req.SenderID)
	if httpClient == nil {
		code = ptl.CodeErrInternal

		return
	}

	senderReq := &sharepkg.SenderTextMessage{
		ReceiverType: req.ReceiverType,
		Receiver:     req.Receiver,
		Text:         req.Text,
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, senderURL,
		bytes.NewReader(senderReq.ToJSONBytes()))
	if err != nil {
		code = ptl.CodeErrInternal
		msg = err.Error()

		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		code = ptl.CodeErrInternal
		msg = err.Error()

		return
	}

	if httpResp == nil {
		code = ptl.CodeErrInternal
		msg = "no http resp"

		return
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		code = ptl.CodeErrInternal
		msg = fmt.Sprintf("http status code: %d", httpResp.StatusCode)

		return
	}

	var resp ptl.ResponseWrapper

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		code = ptl.CodeErrInternal
		msg = err.Error()

		return
	}

	code = resp.Code
	msg = resp.RawMessage

	return
}

func (s *Server) RegisterHandlers(_ *gin.RouterGroup) {

}

func (s *Server) getOrCreateClientForSenderID(senderID sharepkg.SenderID) (cli *http.Client, senderURL, errMsg string) {
	senderURL, ok := s.cfg.Senders[string(senderID)]
	if !ok {
		errMsg = fmt.Sprintf("no url for the %s", senderID)

		return
	}

	i, ok := s.senders.Load(senderID)
	if !ok {
		i, _ = s.senders.LoadOrStore(senderID, &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   2 * time.Second,
					Deadline:  time.Now().Add(3 * time.Second),
					KeepAlive: 2 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 2 * time.Second,
			},
			Timeout: 5 * time.Second,
		})

		ok = true
	}

	if !ok {
		errMsg = fmt.Sprintf("logic error: cant cache http client for %s", senderID)

		return
	}

	cli, ok = i.(*http.Client)
	if !ok {
		errMsg = "cached instance not http client"

		s.senders.Delete(senderID)

		return
	}

	return
}
