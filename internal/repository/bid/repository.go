package bid

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	model2 "tender-service/internal/model"
	"tender-service/internal/model/entity/bid"
	"tender-service/internal/repository/bid/model"
	"tender-service/internal/util"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	bidTableName           = "bid"
	versionTableName       = "bid_version"
	idColumnName           = "id"
	bidIdColumnName        = "bid_id"
	nameColumnName         = "name"
	descriptionColumnName  = "description"
	statusColumnName       = "status"
	tenderIdColumnName     = "tender_id"
	authorTypeColumnName   = "author_type"
	AuthorIdColumnName     = "author_id"
	versionColumnName      = "version"
	bidVersionIdColumnName = "bid_version_id"
	decisionColumnName     = "decision"
	returningAllSuffix     = "RETURNING *"
	bidAndVersionJoin      = "bid_version ON bid.bid_version_id = bid_version.id"
	selectBidSum           = "bid.id, bid_version.name, bid_version.description, bid.status, bid.tender_id, bid.author_type, bid.author_id, bid_version.version, bid.created_at, bid.decision"
)

func NewBidRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SaveBid(ctx context.Context, b bid.Bid) (bid.Bid, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback(ctx)

	bidBuilder := squirrel.Insert(bidTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(statusColumnName, tenderIdColumnName, AuthorIdColumnName, authorTypeColumnName, decisionColumnName).
		Values(b.Status, b.TenderId.String(), b.AuthorId.String(), b.AuthorType, bid.None).
		Suffix(returningAllSuffix)

	sql, args, err := bidBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	savedBid, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Bid])
	if err != nil {
		fmt.Println("x1")
		return bid.Bid{}, err
	}

	versionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(bidIdColumnName, nameColumnName, descriptionColumnName, versionColumnName).
		Values(savedBid.Id.String(), b.Name, b.Description, 1).
		Suffix(returningAllSuffix)

	sql, args, err = versionBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	savedVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.BidVersion])
	if err != nil {
		return bid.Bid{}, err
	}

	setBidVersion := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).
		Set(bidVersionIdColumnName, savedVersion.Id.String()).
		Where(squirrel.Eq{idColumnName: savedBid.Id.String()})

	sql, args, err = setBidVersion.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	err = tx.Commit(ctx)
	if err != nil {
		return bid.Bid{}, err
	}

	return model.MergeBidAndVersionToBid(savedVersion, savedBid), nil
}

func (r *repository) GetBidById(ctx context.Context, id uuid.UUID) (bid.Bid, error) {
	op := "bid_repository.get_by_id"

	builder := squirrel.Select(selectBidSum).PlaceholderFormat(squirrel.Dollar).
		From(bidTableName).Join(bidAndVersionJoin).
		Where(squirrel.Eq{bidIdColumnName: id.String()})

	sql, args, err := builder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	sum, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.BidSum])
	if err != nil {
		return bid.Bid{}, model2.NewNotFoundError(op, err)
	}

	return model.BidSumToBid(sum), nil
}

func (r *repository) GetBidList(ctx context.Context, page util.Page, tenderId uuid.UUID, userId uuid.UUID) ([]bid.Bid, error) {
	builder := squirrel.Select(selectBidSum).PlaceholderFormat(squirrel.Dollar).
		From(bidTableName).Join(bidAndVersionJoin).Offset(uint64(page.Offset)).Limit(uint64(page.Limit))

	if tenderId != uuid.Nil {
		builder = builder.Where(squirrel.Eq{tenderIdColumnName: tenderId.String()})
	}

	log.Println(userId.String())

	if userId != uuid.Nil {
		builder = builder.Where(squirrel.Eq{"bid" + "." + AuthorIdColumnName: userId.String()})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	log.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	sums, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.BidSum])
	if err != nil {
		return nil, err
	}

	return model.BidSumListToBidList(sums), nil
}

func (r *repository) UpdateBidDecision(ctx context.Context, id uuid.UUID, dec bid.Decision) (bid.Bid, error) {
	updateBuilder := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).Set(decisionColumnName, dec).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	return r.GetBidById(ctx, id)
}

func (r *repository) UpdateBidStatus(ctx context.Context, id uuid.UUID, stat bid.Status) (bid.Bid, error) {
	updateBuilder := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).Set(statusColumnName, stat).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		fmt.Println()
		return bid.Bid{}, err
	}

	log.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	return r.GetBidById(ctx, id)
}

func (r *repository) UpdateBid(ctx context.Context, id uuid.UUID, name, description string) (bid.Bid, error) {
	oldVersion, err := r.GetBidById(ctx, id)
	if err != nil {
		return bid.Bid{}, err
	}

	setMap := make(map[string]interface{})

	setMap[versionColumnName] = oldVersion.Version + 1
	setMap[bidIdColumnName] = oldVersion.Id.String()
	setMap[nameColumnName] = oldVersion.Name
	setMap[descriptionColumnName] = oldVersion.Description

	if name != "" {
		setMap[nameColumnName] = name
	}

	if description != "" {
		setMap[descriptionColumnName] = description
	}

	newVersionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		SetMap(setMap).
		Suffix(returningAllSuffix)

	sql, args, err := newVersionBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	newVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.BidVersion])
	if err != nil {
		return bid.Bid{}, err
	}

	updateTenderVersionIdBuilder := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).
		Set(bidVersionIdColumnName, newVersion.Id.String()).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err = updateTenderVersionIdBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql:" + sql)

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	oldVersion.Version = newVersion.Version
	oldVersion.Name = newVersion.Name
	oldVersion.Description = newVersion.Description

	return oldVersion, nil
}

func (r *repository) UpdateTenderStatus(ctx context.Context, id uuid.UUID, stat bid.Status) (bid.Bid, error) {
	updateBuilder := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).Set(statusColumnName, stat).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		fmt.Println()
		return bid.Bid{}, err
	}

	log.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	return r.GetBidById(ctx, id)
}

func (r *repository) RollbackBid(ctx context.Context, id uuid.UUID, ver int) (bid.Bid, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback(ctx)

	curBid, err := r.GetBidById(ctx, id)
	if err != nil {
		return bid.Bid{}, err
	}

	getOldVersionBuilder := squirrel.Select("*").PlaceholderFormat(squirrel.Dollar).
		From(versionTableName).
		Where(squirrel.And{squirrel.Eq{bidIdColumnName: id.String()}, squirrel.Eq{versionColumnName: ver}})

	sql, args, err := getOldVersionBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql1: " + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	oldVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.BidVersion])
	if err != nil {
		return bid.Bid{}, err
	}

	versionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(bidIdColumnName, nameColumnName, descriptionColumnName, versionColumnName).
		Values(curBid.Id.String(), oldVersion.Name, oldVersion.Description, curBid.Version+1).
		Suffix(returningAllSuffix)

	sql, args, err = versionBuilder.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql2: " + sql)

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	savedVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.BidVersion])
	if err != nil {
		return bid.Bid{}, err
	}

	setBidVersion := squirrel.Update(bidTableName).PlaceholderFormat(squirrel.Dollar).
		Set(bidVersionIdColumnName, savedVersion.Id.String()).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err = setBidVersion.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}

	log.Println("sql3: " + sql)

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return bid.Bid{}, err
	}

	rows.Close()

	curBid.Name = oldVersion.Name
	curBid.Description = oldVersion.Description
	curBid.Version += 1

	err = tx.Commit(ctx)
	if err != nil {
		return bid.Bid{}, err
	}

	return curBid, nil
}
