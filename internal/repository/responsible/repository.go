package responsible

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	tableName                = "organization_responsible "
	employeeTableName        = "employee"
	organizationIdColumnName = "organization_id"
	usernameColumnName       = "employee.username"
	employeeUserIdColumnName = "employee.id"
	userIdColumnName         = "organization_responsible.user_id"
)

func NewOrganizationResponsibleRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) UsersHasSimilarOrganization(ctx context.Context, userId uuid.UUID, username string) (bool, error) {
	// squirrel does not support nested queries :_(
	sql := "SELECT EXISTS ( " +
		"SELECT 1 FROM organization_responsible  JOIN employee ON employee.id = organization_responsible.user_id WHERE " +
		"(organization_responsible.organization_id IN ( SELECT organization_id FROM organization_responsible WHERE user_id = $1) AND employee.username = $2))"

	args := []interface{}{userId.String(), username}

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

func (r *repository) IsResponsibleInOrganization(ctx context.Context, username string, organizationId uuid.UUID) (bool, error) {
	builder := squirrel.Select("1").PlaceholderFormat(squirrel.Dollar).
		Prefix("SELECT EXISTS (").From(tableName).Join(employeeTableName + " ON employee.id = organization_responsible.user_id").
		Where(squirrel.And{
			squirrel.Eq{usernameColumnName: username},
			squirrel.Eq{tableName + "." + organizationIdColumnName: organizationId},
		}).Suffix(")")

	sql, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	//fmt.Println("sql:" + sql)

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

func (r *repository) CountEmployeesInOrganization(ctx context.Context, organizationId uuid.UUID) (int, error) {
	builder := squirrel.Select("COUNT (*)").PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Where(squirrel.Eq{organizationIdColumnName: organizationId})

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	rows.Next()

	var count int
	if err = rows.Scan(&count); err != nil {
		return 0, err
	}

	rows.Close()

	return count, nil
}
