package bootstrap

import (
	"context"
	"fmt"

	recordspg "github.com/zchelalo/neuraclinic-records/internal/adapters/postgres"
	appointmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/adapters/grpc"
	appointmentsapp "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application"
	filemanagementadapter "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/filemanagement"
	attachmentsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/adapters/grpc"
	attachmentsapp "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
	familyogramgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/adapters/grpc"
	familyogramapp "github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/application"
	notesgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/notes/adapters/grpc"
	notesapp "github.com/zchelalo/neuraclinic-records/internal/modules/notes/application"
	patientsgrpc "github.com/zchelalo/neuraclinic-records/internal/modules/patients/adapters/grpc"
	patientsapp "github.com/zchelalo/neuraclinic-records/internal/modules/patients/application"
	grpcserver "github.com/zchelalo/neuraclinic-records/internal/server/grpc"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	"go.uber.org/zap"
)

type App struct {
	Server  *grpcserver.Server
	Cleanup func(context.Context) error
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
	familyogramApp := familyogramapp.NewService(repo)
	attachmentsApp := attachmentsapp.NewService(appCfg, repo, filesClient)

	server, err := grpcserver.New(grpcserver.Config{
		Port:            cfg.Port,
		ServiceName:     cfg.ServiceName,
		TLSCertFilePath: cfg.GRPCTLSCertPath,
		TLSKeyFilePath:  cfg.GRPCTLSKeyPath,
	}, logger, grpcserver.Services{
		Patient:     patientsgrpc.NewPatientService(patientsApp),
		Appointment: appointmentsgrpc.NewAppointmentService(appointmentsApp),
		Note:        notesgrpc.NewNoteService(notesApp),
		Familyogram: familyogramgrpc.NewFamilyogramService(familyogramApp),
		Attachment:  attachmentsgrpc.NewAttachmentService(attachmentsApp),
	})
	if err != nil {
		_ = filesClient.Close()
		db.Close()
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	return &App{
		Server: server,
		Cleanup: func(context.Context) error {
			server.GracefulStop()
			_ = filesClient.Close()
			db.Close()
			return nil
		},
	}, nil
}
