package feedback

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tender-service/internal/model/entity"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	tableName                = "feedback"
	idColumnName             = "id"
	bidIdColumnName          = "bid_id"
	descriptionColumnName    = "description"
	usernameColumnName       = "username"
	createdAtColumnName      = "created_at"
	bidTableName             = "bid"
	bidTableIdColumnName     = "bid.id"
	bidTableTenderNameColumn = "bid.tender_id"
)

func NewFeedbackRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SaveFeedback(ctx context.Context, feedback entity.Feedback) (entity.Feedback, error) {
	builder := squirrel.Insert(tableName).PlaceholderFormat(squirrel.Dollar).
		Columns(bidIdColumnName, descriptionColumnName, usernameColumnName).
		Values(feedback.BidId.String(), feedback.Description, feedback.Username).
		Suffix("RETURNING *")

	sql, args, err := builder.ToSql()
	if err != nil {
		return entity.Feedback{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Feedback{}, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Feedback])
	if err != nil {
		return entity.Feedback{}, err
	}

	return result, nil
}

func (r *repository) GetFeedbackListForGroup(ctx context.Context, tenderId uuid.UUID, userId uuid.UUID) ([]entity.Feedback, error) {
	builder := squirrel.Select("feedback.id, feedback.bid_id, feedback.description, feedback.username, feedback.created_at").PlaceholderFormat(squirrel.Dollar).
		From(tableName).Join(bidTableName + " ON bid.id = feedback.bid_id").
		Where(squirrel.And{
			squirrel.Eq{bidTableTenderNameColumn: tenderId.String()},
			squirrel.Eq{"bid.author_id": userId},
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Feedback])
	if err != nil {
		return nil, err
	}

	return result, nil
}
