package models

import (
	"REST-API/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=500"`
	Location    string    `json:"location" validate:"required,min=3,max=100"`
	DateTime    time.Time `json:"dateTime" validate:"required,future_date"`
	UserID      int       `json:"userId"`
}

func (e *Event) Save(ctx context.Context) error {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.DB.ExecContext(ctx, query, e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while saving event")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while saving event")
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id)
	return nil
}

func GetAllEvents(ctx context.Context, page, limit int) ([]Event, int, error) {

	var total int
	countQuery := `SELECT COUNT(*) FROM events`

	err := db.DB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, 0, errors.New("request timeout while counting events")
		}
		return nil, 0, err
	}

	offset := (page - 1) * limit

	query := `SELECT id, name, description, location, dateTime, user_id FROM events LIMIT ? OFFSET ?`
	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, 0, errors.New("request timeout while fetching events")
		}
		return nil, 0, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.DateTime,
			&event.UserID,
		)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func GetEventByID(ctx context.Context, id int) (*Event, error) {
	query := `SELECT id, name, description, location, dateTime, user_id FROM events WHERE id = ?`

	row := db.DB.QueryRowContext(ctx, query, id)

	var event Event
	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("request timeout while fetching event")
		}
		return nil, err
	}

	return &event, nil
}

func (event Event) Update(ctx context.Context) error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	result, err := db.DB.ExecContext(ctx, query, event.Name, event.Description, event.Location, event.DateTime, event.ID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while updating event")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while updating event")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event with id %d not found", event.ID)
	}

	return nil
}

func (event Event) Delete(ctx context.Context) error {
	query := `DELETE FROM events WHERE id = ?`

	result, err := db.DB.ExecContext(ctx, query, event.ID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while deleting event")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while deleting event")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event with id %d not found", event.ID)
	}

	return nil
}
