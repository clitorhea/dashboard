package db

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Template struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Compose     string    `json:"compose"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserCount returns the number of registered users.
func (d *DB) UserCount() (int, error) {
	var count int
	err := d.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// CreateUser inserts a new user with a bcrypt-hashed password.
func (d *DB) CreateUser(username, hashedPassword string) (*User, error) {
	result, err := d.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username, hashedPassword,
	)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &User{ID: id, Username: username, CreatedAt: time.Now()}, nil
}

// GetUserByUsername looks up a user by username.
func (d *DB) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := d.QueryRow(
		"SELECT id, username, password, created_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// CreateSession creates a new session token for a user.
func (d *DB) CreateSession(userID int64) (*Session, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	s := &Session{
		ID:        hex.EncodeToString(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	_, err := d.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		s.ID, s.UserID, s.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetSession validates a session token and returns the associated user ID.
func (d *DB) GetSession(token string) (*Session, error) {
	s := &Session{}
	err := d.QueryRow(
		"SELECT id, user_id, expires_at FROM sessions WHERE id = ? AND expires_at > datetime('now')",
		token,
	).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSession removes a session.
func (d *DB) DeleteSession(token string) error {
	_, err := d.Exec("DELETE FROM sessions WHERE id = ?", token)
	return err
}

// GetUserByID looks up a user by ID.
func (d *DB) GetUserByID(id int64) (*User, error) {
	u := &User{}
	err := d.QueryRow(
		"SELECT id, username, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ListTemplates returns all service templates.
func (d *DB) ListTemplates() ([]Template, error) {
	rows, err := d.Query("SELECT id, name, description, category, compose, icon, created_at FROM templates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var t Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.Compose, &t.Icon, &t.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

// GetTemplate returns a single template by ID.
func (d *DB) GetTemplate(id int64) (*Template, error) {
	t := &Template{}
	err := d.QueryRow(
		"SELECT id, name, description, category, compose, icon, created_at FROM templates WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.Compose, &t.Icon, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// SeedTemplates inserts built-in templates if the table is empty.
func (d *DB) SeedTemplates(templates []Template) error {
	count := 0
	d.QueryRow("SELECT COUNT(*) FROM templates").Scan(&count)
	if count > 0 {
		return nil
	}

	for _, t := range templates {
		_, err := d.Exec(
			"INSERT INTO templates (name, description, category, compose, icon) VALUES (?, ?, ?, ?, ?)",
			t.Name, t.Description, t.Category, t.Compose, t.Icon,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
