package repository

import "github.com/yusufelyldrm/reservation/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res *models.Reservation) error
}
