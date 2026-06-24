package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application/createnote"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application/deletenote"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application/findnote"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application/listnotes"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application/updatenote"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Config = appshared.Config
type Runtime = appshared.Runtime

type Service struct {
	createNote *createnote.UseCase
	listNotes  *listnotes.UseCase
	findNote   *findnote.UseCase
	updateNote *updatenote.UseCase
	deleteNote *deletenote.UseCase
}

func NewService(cfg Config, repo ports.Repository) *Service {
	return NewServiceWithRuntime(cfg, repo, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(cfg Config, repo ports.Repository, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		createNote: createnote.New(repo, runtime),
		listNotes:  listnotes.New(cfg, repo),
		findNote:   findnote.New(repo),
		updateNote: updatenote.New(repo, runtime),
		deleteNote: deletenote.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) CreateNote(ctx context.Context, cmd createnote.Command) (domain.Note, error) {
	return s.createNote.Execute(ctx, cmd)
}

func (s *Service) ListNotes(ctx context.Context, cmd listnotes.Command) (domain.NoteList, error) {
	return s.listNotes.Execute(ctx, cmd)
}

func (s *Service) FindNote(ctx context.Context, psychologistID, id uuid.UUID) (domain.Note, error) {
	return s.findNote.Execute(ctx, findnote.Command{PsychologistID: psychologistID, ID: id})
}

func (s *Service) UpdateNote(ctx context.Context, cmd updatenote.Command) (domain.Note, error) {
	return s.updateNote.Execute(ctx, cmd)
}

func (s *Service) DeleteNote(ctx context.Context, psychologistID, id uuid.UUID) error {
	return s.deleteNote.Execute(ctx, deletenote.Command{PsychologistID: psychologistID, ID: id})
}

type NoteCreateCommand = createnote.Command
type NoteListCommand = listnotes.Command
type NoteUpdateCommand = updatenote.Command
