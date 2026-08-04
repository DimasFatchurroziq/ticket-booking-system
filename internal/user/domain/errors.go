package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("pengguna tidak ditemukan")
	ErrEmailAlreadyTaken  = errors.New("email sudah terdaftar, silakan gunakan email lain")
	ErrInvalidPassword    = errors.New("password yang dimasukkan salah")
	ErrEmailExists        = errors.New("email sudah ada dalam sistem")
	ErrInvalidEmailFormat = errors.New("format email tidak valid")
	ErrInvalidPhoneNumber = errors.New("nomor telepon tidak valid")
)
