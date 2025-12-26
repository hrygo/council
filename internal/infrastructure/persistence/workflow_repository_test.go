package persistence

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/pashagolub/pgxmock/v3"
)

func TestWorkflowRepository_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWorkflowRepository(mock)
	id := uuid.New().String()
	graph := workflow.GraphDefinition{ID: id, Name: "Workflow 1"}
	graphJSON, _ := json.Marshal(graph)

	mock.ExpectQuery("SELECT graph_definition FROM workflows WHERE workflow_uuid = \\$1").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"graph_definition"}).AddRow(graphJSON))

	g, err := repo.Get(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if g.ID != id {
		t.Errorf("expected id %s, got %s", id, g.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWorkflowRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWorkflowRepository(mock)
	id := uuid.New().String()
	graph := &workflow.GraphDefinition{ID: id, Name: "New Workflow"}

	mock.ExpectExec("INSERT INTO workflows").
		WithArgs(id, graph.Name, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), graph)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWorkflowRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWorkflowRepository(mock)
	id := uuid.New().String()
	graph := &workflow.GraphDefinition{ID: id, Name: "Updated Name"}

	mock.ExpectExec("UPDATE workflows").
		WithArgs(id, graph.Name, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), graph)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestWorkflowRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWorkflowRepository(mock)

	mock.ExpectQuery("SELECT workflow_uuid, name, graph_definition, created_at, updated_at FROM workflows").
		WillReturnRows(pgxmock.NewRows([]string{"workflow_uuid", "name", "graph_definition", "created_at", "updated_at"}).
			AddRow("1", "W1", []byte("{}"), time.Now(), time.Now()))

	list, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(list))
	}
}
