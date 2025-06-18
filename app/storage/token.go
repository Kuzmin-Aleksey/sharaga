package storage

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type TokenStorage struct {
	Path string
}

func NewTokenStorage(path string) *TokenStorage {
	return &TokenStorage{
		Path: path,
	}
}

type token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (s *TokenStorage) LoadToken() (string, time.Time) {
	f, err := os.OpenFile(s.Path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("opening token file:", err)
		return "", time.Time{}
	}
	defer f.Close()

	t := new(token)

	if err := json.NewDecoder(f).Decode(t); err != nil {
		log.Println("decoding token:", err)
		return "", time.Time{}
	}

	return t.Token, t.ExpiredAt
}

func (s *TokenStorage) SaveToken(t string, expiredAt time.Time) {
	f, err := os.OpenFile(s.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("opening token file:", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(token{
		Token:     t,
		ExpiredAt: expiredAt,
	}); err != nil {
		log.Println("encoding token:", err)
	}
}

func (s *TokenStorage) DeleteToken() {
	if err := os.Remove(s.Path); err != nil {
		log.Println("deleting token file:", err)
	}
}
