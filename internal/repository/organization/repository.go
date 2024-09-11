package organization

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tender-service/internal/model/entity/organization"
	"tender-service/internal/repository/organization/model"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	tableName    = "organization"
	idColumnName = "id"
)

func NewOrganizationRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) GetOrganizationById(ctx context.Context, organizationId uuid.UUID) (organization.Organization, error) {
	builder := squirrel.Select("*").PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Where(squirrel.Eq{idColumnName: organizationId.String()})

	sql, args, err := builder.ToSql()
	if err != nil {
		return organization.Organization{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return organization.Organization{}, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Organization])
	if err != nil {
		return organization.Organization{}, err
	}

	return model.DbOrganizationToOrganization(result), nil
}

func (r *repository) OrganizationExistById(ctx context.Context, organizationId uuid.UUID) (bool, error) {
	builder := squirrel.Select("1").PlaceholderFormat(squirrel.Dollar).
		Prefix("SELECT EXISTS (").From(tableName).Where(squirrel.Eq{idColumnName: organizationId.String()}).Suffix(")")

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

	return result, nil
}
