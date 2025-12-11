package models

import (
    "context"
    "database/sql"
    "regexp"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
)

func setupMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open sqlmock database: %v", err)
    }
    return db, mock
}

func TestCreateEvent(t *testing.T) {
    dbSQL, mock := setupMock(t)
    defer dbSQL.Close()

    repo := NewEventRepository(dbSQL)

    ev := &Event{
        Title: "Test",
        Description: "desc",
        StartTime: time.Now(),
        EndTime: time.Now().Add(time.Hour),
    }

    // Expect query and return id and created_at
    rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow("550e8400-e29b-41d4-a716-446655440000", time.Now())
    mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO events (title, description, start_time, end_time, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at")).WithArgs(ev.Title, ev.Description, ev.StartTime, ev.EndTime).WillReturnRows(rows)

    ctx := context.Background()
    if err := repo.Create(ctx, ev); err != nil {
        t.Fatalf("Create returned error: %v", err)
    }
    if ev.ID == "" {
        t.Fatalf("expected id to be set")
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestGetByID(t *testing.T) {
    dbSQL, mock := setupMock(t)
    defer dbSQL.Close()

    repo := NewEventRepository(dbSQL)

    id := "550e8400-e29b-41d4-a716-446655440000"
    start := time.Now()
    end := start.Add(time.Hour)
    rows := sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time", "created_at"}).AddRow(id, "T", "d", start, end, time.Now())
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = $1")).WithArgs(id).WillReturnRows(rows)

    ctx := context.Background()
    ev, err := repo.GetByID(ctx, id)
    if err != nil {
        t.Fatalf("GetByID returned error: %v", err)
    }
    if ev.ID != id {
        t.Fatalf("expected id %s, got %s", id, ev.ID)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestGetAll(t *testing.T) {
    dbSQL, mock := setupMock(t)
    defer dbSQL.Close()

    repo := NewEventRepository(dbSQL)

    start := time.Now()
    end := start.Add(time.Hour)
    rows := sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time", "created_at"}).
        AddRow("id1", "T1", "d1", start, end, time.Now()).
        AddRow("id2", "T2", "d2", start, end, time.Now())

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, description, start_time, end_time, created_at FROM events")).WillReturnRows(rows)

    ctx := context.Background()
    list, err := repo.GetAll(ctx)
    if err != nil {
        t.Fatalf("GetAll returned error: %v", err)
    }
    if len(list) != 2 {
        t.Fatalf("expected 2 events, got %d", len(list))
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestUpdateEvent(t *testing.T) {
    dbSQL, mock := setupMock(t)
    defer dbSQL.Close()

    repo := NewEventRepository(dbSQL)
    ev := &Event{ID: "id1", Title: "T1", Description: "d1", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)}

    mock.ExpectExec(regexp.QuoteMeta("UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4 WHERE id = $5")).WithArgs(ev.Title, ev.Description, ev.StartTime, ev.EndTime, ev.ID).WillReturnResult(sqlmock.NewResult(0, 1))

    ctx := context.Background()
    if err := repo.Update(ctx, ev); err != nil {
        t.Fatalf("Update returned error: %v", err)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestDeleteEvent(t *testing.T) {
    dbSQL, mock := setupMock(t)
    defer dbSQL.Close()

    repo := NewEventRepository(dbSQL)
    id := "id1"
    mock.ExpectExec(regexp.QuoteMeta("DELETE FROM events WHERE id = $1")).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

    ctx := context.Background()
    if err := repo.Delete(ctx, id); err != nil {
        t.Fatalf("Delete returned error: %v", err)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}
