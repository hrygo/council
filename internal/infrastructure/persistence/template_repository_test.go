package persistence

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/pashagolub/pgxmock/v3"
)

func TestTemplateRepository_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewTemplateRepository(mock)
	id := "template-1"
	tpl := workflow.Template{ID: id, Name: "Template 1", Category: "custom"}
	graphJSON, _ := json.Marshal(tpl.Graph)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS templates").WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mock.ExpectQuery("SELECT id, name, description, category, is_system, graph, created_at, updated_at FROM templates WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "category", "is_system", "graph", "created_at", "updated_at"}).
			AddRow(id, "Template 1", "desc", "custom", false, graphJSON, time.Now(), time.Now()))

	t1, err := repo.Get(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if t1.ID != id {
		t.Errorf("expected id %s, got %s", id, t1.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTemplateRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewTemplateRepository(mock)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS templates").WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mock.ExpectQuery("SELECT id, name, description, category, is_system, graph, created_at, updated_at FROM templates").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "category", "is_system", "graph", "created_at", "updated_at"}).
			AddRow("1", "T1", "", "custom", false, []byte("{}"), time.Now(), time.Now()))

	list, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 template, got %d", len(list))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTemplateRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewTemplateRepository(mock)
	tpl := &workflow.Template{
		ID:       "tpl-1",
		Name:     "New Template",
		Category: "custom",
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS templates").WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mock.ExpectExec("INSERT INTO templates").
		WithArgs(tpl.ID, tpl.Name, tpl.Description, string(tpl.Category), tpl.IsSystem, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), tpl)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTemplateRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewTemplateRepository(mock)
	id := "tpl-1"

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS templates").WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mock.ExpectExec("DELETE FROM templates").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
