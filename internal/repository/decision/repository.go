package decision

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tender-service/internal/model/entity/decision"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	tableName          = "decision"
	verdictColumnName  = "verdict"
	usernameColumnName = "username"
	bidIdColumnName    = "bid_id"
)

func NewDecisionRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SaveDecision(ctx context.Context, dec decision.Decision) (decision.Decision, error) {
	builder := squirrel.Insert(tableName).PlaceholderFormat(squirrel.Dollar).
		Columns(bidIdColumnName, verdictColumnName, usernameColumnName).
		Values(dec.BidId, dec.Verdict, dec.Username).
		Suffix("RETURNING *")

	sql, args, err := builder.ToSql()
	if err != nil {
		return decision.Decision{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return decision.Decision{}, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[decision.Decision])
	if err != nil {
		return decision.Decision{}, err
	}

	return result, nil
}

func (r *repository) CountDecisionForBid(ctx context.Context, bidId uuid.UUID) (int, error) {
	builder := squirrel.Select("COUNT(*)").PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Where(squirrel.Eq{bidIdColumnName: bidId.String()})

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
