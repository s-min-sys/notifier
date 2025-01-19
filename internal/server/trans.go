package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/s-min-sys/notifier-share/pkg/model"
	"github.com/sgostarter/libeasygo/ptl"
)

func (s *Server) trans(senderBy model.SenderBy, urlPath string, body []byte, rResp any) (code ptl.Code, msg string) {
	if senderBy == model.SenderByAll {
		code = ptl.CodeErrInvalidArgs

		return
	}

	cli, u, msg := s.getOrCreateClientForSenderID(senderBy)
	if msg != "" {
		code = ptl.CodeErrInternal

		return
	}

	rURL, err := url.JoinPath(u, urlPath)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, rURL,
		bytes.NewReader(body))
	if err != nil {
		code = ptl.CodeErrInternal
		msg = err.Error()

		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := cli.Do(httpReq)
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

	resp.Resp = rResp

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
