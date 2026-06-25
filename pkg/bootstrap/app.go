package bootstrap

import (
	"context"
	"fmt"

	recordspg "github.com/zchelalo/neuraclinic-records/internal/adapters/postgres"
	appointmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/adapters/grpc"
	appointmentsapp "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application"
	filemanagementadapter "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/filemanagement"
	attachmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/grpc"
	attachmentsrabbit "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/rabbitmq"
	attachmentsapp "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
	familiogramgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/adapters/grpc"
	familiogramapp "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/application"
	notesgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/notes/adapters/grpc"
	notesapp "github.com/zchelalo/neuraclinic-records/internal/modules/notes/application"
	patientsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/patients/adapters/grpc"
	patientsapp "github.com/zchelalo/neuraclinic-records/internal/modules/patients/application"
	grpcserver "github.com/zchelalo/neuraclinic-records/internal/server/grpc"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	"go.uber.org/zap"
)

type App struct {
	Server   *grpcserver.Server
	Consumer *attachmentsrabbit.Consumer
	Cleanup  func(context.Context) error
}

func InitApp(ctx context.Context, logger *zap.Logger, cfg Config) (*App, error) {
	db, err := NewDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db: %w", err)
	}

	filesClient, err := filemanagementadapter.New(filemanagementadapter.Config{
		Addr:               cfg.FileManagementGRPCAddr,
		TLSEnabled:         cfg.FileManagementGRPCTLSEnabled,
		CACertPath:         cfg.FileManagementGRPCCACertPath,
		InsecureSkipVerify: cfg.FileManagementGRPCInsecureSkipVerify,
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot initialize file-management client: %w", err)
	}

	repo := recordspg.NewRepository(db)
	appCfg := appshared.Config{
		PaginationLimitDefault: cfg.PaginationLimitDefault,
		PaginationLimitMax:     cfg.PaginationLimitMax,
	}
	patientsApp := patientsapp.NewService(appCfg, repo)
	appointmentsApp := appointmentsapp.NewService(appCfg, repo)
	notesApp := notesapp.NewService(appCfg, repo)
	familiogramApp := familiogramapp.NewService(repo)
	attachmentsApp := attachmentsapp.NewService(appCfg, repo, filesClient)
	consumer, err := attachmentsrabbit.NewConsumer(attachmentsrabbit.Config{
		URL:        cfg.RabbitMQURL,
		Exchange:   cfg.RabbitMQExchange,
		Queue:      cfg.RabbitMQQueue,
		RoutingKey: cfg.RabbitMQRoutingKey,
		DLX:        cfg.RabbitMQDLX,
		DLQ:        cfg.RabbitMQDLQ,
		Prefetch:   cfg.RabbitMQPrefetch,
	}, attachmentsrabbit.NewHandler(attachmentsApp), logger)
	if err != nil {
		_ = filesClient.Close()
		db.Close()
		return nil, fmt.Errorf("cannot initialize rabbitmq consumer: %w", err)
	}

	server, err := grpcserver.New(grpcserver.Config{
		Port:            cfg.Port,
		ServiceName:     cfg.ServiceName,
		TLSCertFilePath: cfg.GRPCTLSCertPath,
		TLSKeyFilePath:  cfg.GRPCTLSKeyPath,
	}, logger, grpcserver.Services{
		Patient:     patientsgrpc.NewPatientService(patientsApp),
		Appointment: appointmentsgrpc.NewAppointmentService(appointmentsApp),
		Note:        notesgrpc.NewNoteService(notesApp),
		Familiogram: familiogramgrpc.NewFamiliogramService(familiogramApp),
		Attachment:  attachmentsgrpc.NewAttachmentService(attachmentsApp),
	})
	if err != nil {
		_ = consumer.Close()
		_ = filesClient.Close()
		db.Close()
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	return &App{
		Server:   server,
		Consumer: consumer,
		Cleanup: func(context.Context) error {
			server.GracefulStop()
			_ = consumer.Close()
			_ = filesClient.Close()
			db.Close()
			return nil
		},
	}, nil
}
