package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/controller"
	"github.com/sbasestarter/proto-repo/gen/protorepo-post-go"
)

type PostServer struct {
	controller *controller.Controller
}

func NewPostServer(ctx context.Context, cfg *config.Config) *PostServer {
	return &PostServer{
		controller: controller.NewController(ctx, cfg),
	}
}

func (s *PostServer) makeStatus(status postpb.PostStatus, err error) *postpb.ServerStatus {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return s.makeStatusWithMsg(status, msg)
}

func (s *PostServer) makeStatusWithMsg(status postpb.PostStatus, msg string) *postpb.ServerStatus {
	if msg == "" {
		msg = status.String()
	}
	return &postpb.ServerStatus{
		Status: status,
		Msg:    msg,
	}
}

func (s *PostServer) SendTemplate(ctx context.Context, req *postpb.SendTemplateRequest) (*postpb.SendTemplateResponse, error) {
	err := s.controller.SendTemplate(ctx, req.To, req.ProtocolType, req.TemplateId, req.Vars)
	if err != nil {
		return &postpb.SendTemplateResponse{
			Status: s.makeStatus(postpb.PostStatus_PS_FAILED, err),
		}, nil
	}
	return &postpb.SendTemplateResponse{
		Status: s.makeStatus(postpb.PostStatus_PS_SUCCESS, nil),
	}, nil
}

func (s *PostServer) HTTPRegister(r *mux.Router) {
	r.HandleFunc("/send", s.httpSendTemplate)
}

func (s *PostServer) httpSendTemplate(w http.ResponseWriter, r *http.Request) {
	protocolType := r.URL.Query().Get("protocol_type")

	to := r.URL.Query()["to"]
	templateID := r.URL.Query().Get("template_id")
	vars := r.URL.Query()["vars"]
	err := s.controller.SendTemplate(r.Context(), to, protocolType, templateID, vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}
