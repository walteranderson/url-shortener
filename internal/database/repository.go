package database

import (
	"database/sql"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Link struct {
	Url  string
	Code string
}

func (r *Repository) CreateLink(url string) (*Link, error) {
	code, err := r.newCode()
	if err != nil {
		return nil, err
	}
	query := "INSERT INTO links (url, code) VALUES ($1, $2)"
	if _, err := r.db.Exec(query, url, code); err != nil {
		return nil, err
	}
	link := &Link{Url: url, Code: code}
	return link, nil
}

func (r *Repository) GetLink(code string) (*Link, error) {
	var url string
	query := "SELECT url FROM links WHERE code = $1"
	err := r.db.QueryRow(query, code).Scan(&url)
	if err != nil {
		return nil, err
	}
	link := &Link{Url: url, Code: code}
	return link, nil
}

func (r *Repository) newCode() (string, error) {
	v := make([]rune, 6)
	for i := range v {
		v[i] = letters[rand.Intn(len(letters))]
	}
	// TODO add a db query to check that the code is unique
	return string(v), nil
}
