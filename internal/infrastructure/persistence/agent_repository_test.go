package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAgentRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAgentRepository(mock)
	id := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "name", "avatar", "description", "persona_prompt", "model_config", "capabilities", "created_at", "updated_at"}).
		AddRow(id, "Agent 1", "avatar", "desc", "persona", agent.ModelConfig{}, agent.Capabilities{}, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, avatar, description, persona_prompt, model_config, capabilities, created_at, updated_at FROM agents WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	a, err := repo.GetByID(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if a.ID != id {
		t.Errorf("expected id %v, got %v", id, a.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAgentRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAgentRepository(mock)
	a := &agent.Agent{
		Name: "New Agent",
	}

	mock.ExpectQuery("INSERT INTO agents").
		WithArgs(a.Name, a.Avatar, a.Description, a.PersonaPrompt, a.ModelConfig, a.Capabilities, pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(uuid.New(), time.Now(), time.Now()))

	err = repo.Create(context.Background(), a)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAgentRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAgentRepository(mock)

	mock.ExpectQuery("SELECT id, name, avatar, description, persona_prompt, model_config, capabilities, created_at, updated_at FROM agents").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "avatar", "description", "persona_prompt", "model_config", "capabilities", "created_at", "updated_at"}).
			AddRow(uuid.New(), "A1", "", "", "", agent.ModelConfig{}, agent.Capabilities{}, time.Now(), time.Now()))

	list, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 agent, got %d", len(list))
	}
}

func TestAgentRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAgentRepository(mock)
	id := uuid.New()
	a := &agent.Agent{ID: id, Name: "Updated Name"}

	mock.ExpectExec("UPDATE agents").
		WithArgs(a.Name, a.Avatar, a.Description, a.PersonaPrompt, a.ModelConfig, a.Capabilities, pgxmock.AnyArg(), a.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), a)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAgentRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAgentRepository(mock)
	id := uuid.New()

	mock.ExpectExec("DELETE FROM agents").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
