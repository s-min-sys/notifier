package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/protocx/pjson"
	"github.com/s-min-sys/notifier-share/v2/pkg/model"
)

func (s *Server) trans(senderBy model.SenderBy, urlPath string, body []byte, rResp any) (code errorx.Code, msg string) {
	if senderBy == model.SenderByAll {
		code = errorx.CodeErrInvalidArgs

		return
	}

	cli, u, msg := s.getOrCreateClientForSenderID(senderBy)
	if msg != "" {
		code = errorx.CodeErrInternal

		return
	}

	rURL, err := url.JoinPath(u, urlPath)
	if err != nil {
		code = errorx.CodeErrInternal

		msg = err.Error()

		return
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, rURL,
		bytes.NewReader(body))
	if err != nil {
		code = errorx.CodeErrInternal
		msg = err.Error()

		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := cli.Do(httpReq)
	if err != nil {
		code = errorx.CodeErrInternal
		msg = err.Error()

		return
	}

	if httpResp == nil {
		code = errorx.CodeErrInternal
		msg = "no http resp"

		return
	}

	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		code = errorx.CodeErrInternal
		msg = fmt.Sprintf("http status code: %d", httpResp.StatusCode)

		return
	}

	var resp pjson.ResponseWrapper

	resp.Resp = rResp

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		code = errorx.CodeErrInternal
		msg = err.Error()

		return
	}

	code = resp.Code
	msg = resp.RawMessage

	return
}
