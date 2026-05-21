package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/base/logx"
	"github.com/gin-gonic/gin"
	sharepkg "github.com/s-min-sys/notifier-share/v2/pkg"
	"github.com/s-min-sys/notifier-share/v2/pkg/model"
	"github.com/s-min-sys/notifier/v2/internal/config"
)

func NewServer(cfg *config.Config, logger logx.Wrapper) *Server {
	if logger == nil {
		logger = logx.NewNopLoggerWrapper()
	}

	logger = logger.WithFields(logx.StringField(logx.ClsKey, "server"))

	if cfg == nil {
		logger.Fatal("no config")
	}

	s := &Server{
		cfg:    cfg,
		logger: logger,
	}

	s.init()

	return s
}

type Server struct {
	wg     sync.WaitGroup
	cfg    *config.Config
	logger logx.Wrapper

	senders sync.Map

	allSenders []model.SenderBy
}

func (s *Server) init() {
	for sender := range s.cfg.Senders {
		s.allSenders = append(s.allSenders, model.SenderBy(sender))
	}

	s.wg.Add(1)

	go s.httpServerRoutine()
}

func (s *Server) httpServerRoutine() {
	defer s.wg.Done()

	sharepkg.RunCommandServer[model.TextMessage](s.cfg.Listens, "", s, s, s.logger)
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) SendTextMessage(req *model.TextMessage, _ sharepkg.Storage) (code errorx.Code, msg string) {
	if req.SenderBy != model.SenderByAll {
		return s.sendTextMessage(req)
	}

	var errorCount int

	m := make(map[string]string)

	for _, sender := range s.allSenders {
		req.SenderBy = sender

		code, msg = s.sendTextMessage(req)
		m[string(sender)] = fmt.Sprintf("%d: %s", code, msg)

		if code != errorx.CodeSuccess {
			errorCount++
		}
	}

	if errorCount == 0 {
		return errorx.CodeSuccess, ""
	}

	d, _ := json.Marshal(m)

	return errorx.CodeErrInternal, string(d)
}

func (s *Server) sendTextMessage(req *model.TextMessage) (errorx.Code, string) {
	senderReq := &model.TextMessage{
		SendMessageTarget: model.SendMessageTarget{
			SenderBy: req.SenderBy,
			BizCode:  req.BizCode,
			ToType:   req.ToType,
			To:       req.To,
			FindOpts: req.FindOpts,
		},
		Text: req.Text,
	}

	return s.trans(req.SenderBy, sharepkg.URLSendTextMessage, senderReq.ToJSONBytes(), nil)
}

func (s *Server) RegisterHandlers(_ *gin.RouterGroup) {

}

func (s *Server) getOrCreateClientForSenderID(senderID model.SenderBy) (cli *http.Client, senderURL, errMsg string) {
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
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 10 * time.Second,
			},
			Timeout: 5 * time.Second,
		})
	}

	cli, ok = i.(*http.Client)
	if !ok {
		errMsg = "cached instance not http client"

		s.senders.Delete(senderID)

		return
	}

	return
}

//
// Storage
//

func (s *Server) AddAdminUser(req *model.AdminUserAdd) (code errorx.Code, msg string) {
	return s.trans(req.SenderBy, sharepkg.URLAddAdminUser, req.ToJSONBytes(), nil)
}

func (s *Server) GetAdminUsers(req *model.AdminUserGet) (users []model.UserG, code errorx.Code, msg string) {
	code, msg = s.trans(req.SenderBy, sharepkg.URLGetAdminUsers, req.ToJSONBytes(), &users)

	return
}

func (s *Server) FilterSenderTargets(_ model.SendMessageTarget) []model.SenderTarget {
	return nil
}

func (s *Server) FindUser(_ string, _ int) (*model.User, error) {
	return nil, errorx.ErrUnimplemented
}

func (s *Server) AddUser(_ model.User) (err error) {
	err = errorx.ErrUnimplemented

	return
}
