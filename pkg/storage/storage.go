package storage

import "github.com/MikhailLipanin/html-parser/pkg/parsing"

type Storage interface {
	ReadAllContent() []parsing.ErrorType
	IsPresent(string) bool
	GetValById(string) string
	UpdateValById(string, string) error
	AddValById(string, string) error
	DeleteById(string) error
}
