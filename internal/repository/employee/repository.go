package employee

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model2 "tender-service/internal/model"
	"tender-service/internal/model/entity"
	"tender-service/internal/repository/employee/model"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	tableName          = "employee"
	idColumnName       = "id"
	usernameColumnName = "username"
)

var (
	errEmployeeNotFound = fmt.Errorf("employee not found")
)

func NewEmployeeRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) GetEmployeeByUsername(ctx context.Context, username string) (entity.Employee, error) {
	op := "employee_repository.get_employee_by_username"

	builder := squirrel.Select("*").PlaceholderFormat(squirrel.Dollar).
		From(tableName).Where(squirrel.Eq{usernameColumnName: username})

	sql, args, err := builder.ToSql()
	if err != nil {
		return entity.Employee{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Employee{}, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Employee])
	if err != nil {
		return entity.Employee{}, model2.NewNotAuthorizedError(op, errEmployeeNotFound)
	}

	return model.DbEmployeeToEmployee(result), nil
}

func (r *repository) EmployeeExistByUsername(ctx context.Context, username string) (bool, error) {
	builder := squirrel.Select("1").PlaceholderFormat(squirrel.Dollar).
		Prefix("SELECT EXISTS (").From(tableName).Where(squirrel.Eq{usernameColumnName: username}).Suffix(")")

	sql, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return false, err
	}

	rows.Next()

	var result bool
	err = rows.Scan(&result)
	if err != nil {
		return false, err
	}

	rows.Close()

	return result, nil
}

func (r *repository) GetEmployeeById(ctx context.Context, id uuid.UUID) (entity.Employee, error) {
	op := "employee_repository.get_employee_by_id"
	builder := squirrel.Select("*").PlaceholderFormat(squirrel.Dollar).
		From(tableName).Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err := builder.ToSql()
	if err != nil {
		return entity.Employee{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Employee{}, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Employee])
	if err != nil {
		return entity.Employee{}, model2.NewNotFoundError(op, errEmployeeNotFound)
	}

	return model.DbEmployeeToEmployee(result), nil
}

func (r *repository) EmployeeExistById(ctx context.Context, id uuid.UUID) (bool, error) {
	builder := squirrel.Select("1").PlaceholderFormat(squirrel.Dollar).
		Prefix("SELECT EXISTS (").From(tableName).Where(squirrel.Eq{idColumnName: id.String()}).Suffix(")")

	sql, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return false, err
	}

	rows.Next()

	var result bool
	err = rows.Scan(&result)
	if err != nil {
		return false, err
	}

	rows.Close()

	return result, nil
}
