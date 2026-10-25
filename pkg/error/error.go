package apperror

import "errors"

var (
	ErrNotFound         = errors.New("data tidak ditemukan")
	ErrUnauthorized     = errors.New("unauthorized: token tidak valid atau tidak disertakan")
	ErrInvalidLogin     = errors.New("username atau password salah")
	ErrUserExists       = errors.New("username sudah digunakan")
	ErrUserNotFound     = errors.New("user tidak ditemukan")
	ErrInvalidYear      = errors.New("tahun rilis harus di antara tahun 1980 dan 2024")
	ErrValidation       = errors.New("validasi data gagal")
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
	ErrBookNotFound     = errors.New("buku tidak ditemukan")
)
