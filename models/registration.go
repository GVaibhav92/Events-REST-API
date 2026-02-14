package models

import (
	"REST-API/db"
	"errors"
)

type Registration struct {
	ID      int `json:"id"`
	EventID int `json:"eventId"`
	UserID  int `json:"userId"`
}

func (r *Registration) Save() error {
	event, err := GetEventByID(r.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	alreadyRegistered, err := IsUserRegistered(r.EventID, r.UserID)
	if err != nil {
		return err
	}
	if alreadyRegistered {
		return errors.New("already registered for this event")
	}

	query := `INSERT INTO registrations(event_id, user_id) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.EventID, r.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = int(id)
	return nil
}

func (r *Registration) Cancel() error {
	event, err := GetEventByID(r.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	alreadyRegistered, err := IsUserRegistered(r.EventID, r.UserID)
	if err != nil {
		return err
	}
	if !alreadyRegistered {
		return errors.New("you are not registered for this event")
	}

	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.EventID, r.UserID)
	return err
}

func IsUserRegistered(eventID, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?`
	row := db.DB.QueryRow(query, eventID, userID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
