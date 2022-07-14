package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/controller"
	postpb "github.com/sbasestarter/proto-repo/gen/protorepo-post-go"
	sharepb "github.com/sbasestarter/proto-repo/gen/protorepo-share-go"
	"github.com/sgostarter/i/l"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServer struct {
	controller *controller.Controller
}

func NewPostServer(ctx context.Context, cfg *config.Config, logger l.Wrapper) *PostServer {
	return &PostServer{
		controller: controller.NewController(ctx, cfg, logger),
	}
}

func (s *PostServer) SendTemplate(ctx context.Context, req *postpb.SendTemplateRequest) (*sharepb.Empty, error) {
	err := s.controller.SendTemplate(ctx, req.Tos, req.ProtocolType, req.TemplateId, req.Vars)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.Empty{}, nil
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
