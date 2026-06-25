package grpcserver

import (
	"crypto/tls"
	"fmt"
	"net"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	appointmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/adapters/grpc"
	attachmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/grpc"
	familiogramgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/adapters/grpc"
	notesgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/notes/adapters/grpc"
	patientsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/patients/adapters/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Port            int
	ServiceName     string
	TLSCertFilePath string
	TLSKeyFilePath  string
}

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

type Services struct {
	Patient     *patientsgrpc.PatientService
	Appointment *appointmentsgrpc.AppointmentService
	Note        *notesgrpc.NoteService
	Familiogram *familiogramgrpc.FamiliogramService
	Attachment  *attachmentsgrpc.AttachmentService
}

func New(cfg Config, logger *zap.Logger, appServices Services) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(cfg.TLSCertFilePath, cfg.TLSKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("load grpc tls key pair: %w", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(&tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
		})),
		grpc.UnaryInterceptor(UnaryInterceptor(logger, cfg.ServiceName)),
	)

	recordv1.RegisterPatientServiceServer(grpcServer, appServices.Patient)
	recordv1.RegisterAppointmentServiceServer(grpcServer, appServices.Appointment)
	recordv1.RegisterNoteServiceServer(grpcServer, appServices.Note)
	recordv1.RegisterFamiliogramServiceServer(grpcServer, appServices.Familiogram)
	recordv1.RegisterAttachmentServiceServer(grpcServer, appServices.Attachment)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
	}, nil
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
