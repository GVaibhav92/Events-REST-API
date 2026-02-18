package models

import (
	"REST-API/db"
	"database/sql"
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

func (e *Event) Save() error {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)
	`
	//precompiles SQL,prevents SQL injection
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id) //DB id assigned to struct
	return nil
}

func GetAllEvents(page, limit int) ([]Event, int, error) {
	//get total number of events in the table
	countRow := db.DB.QueryRow(`SELECT COUNT(*) FROM events`)

	var total int
	err := countRow.Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	//LIMIT - Maximum number of rows to return
	// OFFSET - Number of rows to skip before starting

	offset := (page - 1) * limit

	query := `SELECT * FROM events LIMIT ? OFFSET ?`
	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
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

func GetEventByID(id int) (*Event, error) {
	query := `SELECT * FROM events WHERE id = ?`

	row := db.DB.QueryRow(query, id)

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
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	if err != nil {
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

func (event Event) Delete() error {
	query := `DELETE FROM events WHERE id = ?`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(event.ID)
	if err != nil {
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
