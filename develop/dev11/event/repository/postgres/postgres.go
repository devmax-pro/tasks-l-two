package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"

	"module-calendar/event"
)

type PostgresEventRepository struct {
	*sql.DB
}

func NewPostgresEventRepository(db *sql.DB) event.EventRepository {
	return &PostgresEventRepository{
		DB: db,
	}
}

func NewPostgresDB(path string) (*sql.DB, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *PostgresEventRepository) Create(user_id uint64, e event.Event) (event.Event, error) {
	query := `INSERT INTO events (user_id, title, date) VALUES ($1, $2, $3) RETURNING id`
	err := p.QueryRow(query, user_id, e.Title, e.Date).Scan(&e.ID)
	if err != nil {
		return event.Event{}, err
	}
	return e, nil
}

func (p *PostgresEventRepository) Update(user_id uint64, e event.Event) error {
	query := `UPDATE events SET title = $1, date = $2 WHERE id = $3 AND user_id = $4`
	_, err := p.Exec(query, e.Title, e.Date, e.ID, user_id)
	return err
}

func (p *PostgresEventRepository) Delete(user_id uint64, event_id uint64) error {
	query := `DELETE FROM events WHERE id = $1 AND user_id = $2`
	_, err := p.Exec(query, event_id, user_id)
	return err
}

func (p *PostgresEventRepository) GetForDay(user_id uint64, day time.Time) ([]event.Event, error) {
	query := `SELECT id, title, date FROM events WHERE user_id = $1 AND date = $2`
	rows, err := p.Query(query, user_id, day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Date); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (p *PostgresEventRepository) GetForWeek(user_id uint64, week time.Time) ([]event.Event, error) {
	// Assuming week starts on Monday
	startOfWeek := week.AddDate(0, 0, -int(week.Weekday()-time.Monday))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	return p.getEventsBetweenDates(user_id, startOfWeek, endOfWeek)
}

func (p *PostgresEventRepository) GetForMonth(user_id uint64, month time.Time) ([]event.Event, error) {
	// Assuming month starts on the first day
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	return p.getEventsBetweenDates(user_id, startOfMonth, endOfMonth)
}

func (p *PostgresEventRepository) getEventsBetweenDates(user_id uint64, start, end time.Time) ([]event.Event, error) {
	query := `SELECT id, title, date FROM events WHERE user_id = $1 AND date >= $2 AND date < $3`
	rows, err := p.Query(query, user_id, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Date); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
