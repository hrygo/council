package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/group"
	"github.com/pashagolub/pgxmock/v3"
)

func TestGroupRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewGroupRepository(mock)
	id := uuid.New()

	rows := pgxmock.NewRows([]string{"group_uuid", "name", "icon", "system_prompt", "default_agent_uuids", "created_at", "updated_at"}).
		AddRow(id, "Group 1", "icon", "prompt", []uuid.UUID{}, time.Now(), time.Now())

	mock.ExpectQuery("SELECT group_uuid, name, icon, system_prompt, default_agent_uuids, created_at, updated_at FROM groups WHERE group_uuid = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	g, err := repo.GetByID(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if g.ID != id {
		t.Errorf("expected id %v, got %v", id, g.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGroupRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewGroupRepository(mock)
	id := uuid.New()
	g := &group.Group{
		Name:              "New Group",
		Icon:              strPtr("icon"),
		SystemPrompt:      strPtr("prompt"),
		DefaultAgentUUIDs: []uuid.UUID{},
	}

	mock.ExpectQuery("INSERT INTO groups").
		WithArgs(g.Name, g.Icon, g.SystemPrompt, g.DefaultAgentUUIDs, pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"group_uuid", "created_at", "updated_at"}).AddRow(id, time.Now(), time.Now()))

	err = repo.Create(context.Background(), g)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if g.ID != id {
		t.Errorf("expected id %v, got %v", id, g.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGroupRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewGroupRepository(mock)

	mock.ExpectQuery("SELECT group_uuid, name, icon, system_prompt, default_agent_uuids, created_at, updated_at FROM groups ORDER BY created_at DESC").
		WillReturnRows(pgxmock.NewRows([]string{"group_uuid", "name", "icon", "system_prompt", "default_agent_uuids", "created_at", "updated_at"}).
			AddRow(uuid.New(), "G1", strPtr("icon"), strPtr("prompt"), []uuid.UUID{}, time.Now(), time.Now()))

	list, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 group, got %d", len(list))
	}
}

func TestGroupRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewGroupRepository(mock)
	id := uuid.New()
	g := &group.Group{ID: id, Name: "Updated Name"}

	mock.ExpectExec("UPDATE groups").
		WithArgs(g.Name, g.Icon, g.SystemPrompt, g.DefaultAgentUUIDs, pgxmock.AnyArg(), g.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), g)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGroupRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewGroupRepository(mock)
	id := uuid.New()

	mock.ExpectExec("DELETE FROM groups").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
