package model

import (
	"database/sql"
	"github.com/google/uuid"
	"tender-service/internal/model/entity"
	"time"
)

type Employee struct {
	Id        uuid.UUID      `db:"id"`
	Username  string         `db:"username"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func DbEmployeeToEmployee(employee Employee) entity.Employee {
	return entity.Employee{
		Id:        employee.Id,
		Username:  employee.Username,
		FirstName: employee.FirstName.String,
		LastName:  employee.LastName.String,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}
}
