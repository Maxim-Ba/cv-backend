package repository

import (
	"database/sql"
	"fmt"
)

// Profile хранит контактные данные владельца CV
type Profile struct {
	ID       int64
	FullName string
	Title    string
	About    *string
	Email    *string
	Telegram *string
	GitHub   *string
	Phone    *string
}

// ProfileRepo репозиторий для чтения профиля
type ProfileRepo struct {
	db *sql.DB
}

// NewProfileRepo создает новый экземпляр ProfileRepo
func NewProfileRepo(db *sql.DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

// Get возвращает единственную запись профиля
func (r *ProfileRepo) Get() (Profile, error) {
	const query = `
		SELECT id, full_name, title, about, email, telegram, github, phone
		FROM profile
		LIMIT 1
	`
	var p Profile
	err := r.db.QueryRow(query).Scan(
		&p.ID,
		&p.FullName,
		&p.Title,
		&p.About,
		&p.Email,
		&p.Telegram,
		&p.GitHub,
		&p.Phone,
	)
	if err == sql.ErrNoRows {
		return Profile{}, fmt.Errorf("profile not found")
	}
	if err != nil {
		return Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	return p, nil
}
