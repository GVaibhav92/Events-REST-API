package models

import (
	"REST-API/db"
	"context"
	"errors"
)

type Registration struct {
	ID      int `json:"id"`
	EventID int `json:"eventId"`
	UserID  int `json:"userId"`
}

func (r *Registration) Save(ctx context.Context) error {
	event, err := GetEventByID(ctx, r.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	alreadyRegistered, err := IsUserRegistered(ctx, r.EventID, r.UserID)
	if err != nil {
		return err
	}
	if alreadyRegistered {
		return errors.New("already registered for this event")
	}

	query := `INSERT INTO registrations(event_id, user_id) VALUES (?, ?)`

	result, err := db.DB.ExecContext(ctx, query, r.EventID, r.UserID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while registering for event")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while registering for event")
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = int(id)
	return nil
}

func (r *Registration) Cancel(ctx context.Context) error {
	event, err := GetEventByID(ctx, r.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	alreadyRegistered, err := IsUserRegistered(ctx, r.EventID, r.UserID)
	if err != nil {
		return err
	}
	if !alreadyRegistered {
		return errors.New("you are not registered for this event")
	}

	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`

	_, err = db.DB.ExecContext(ctx, query, r.EventID, r.UserID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while canceling registration")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while canceling registration")
		}
		return err
	}

	return nil
}

func IsUserRegistered(ctx context.Context, eventID, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?`

	row := db.DB.QueryRowContext(ctx, query, eventID, userID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, errors.New("request timeout while checking registration")
		}
		return false, err
	}

	return count > 0, nil
}
