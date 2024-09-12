package tender

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model2 "tender-service/internal/model"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/repository/tender/model"
	"tender-service/internal/util"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	versionTableName          = "tender_version"
	tenderTableName           = "tender"
	idColumnName              = "id"
	tenderIdColumnName        = "tender_id"
	nameColumnName            = "name"
	descriptionColumnName     = "description"
	statusColumnName          = "status"
	serviceTypeColumnName     = "service_type"
	versionColumnName         = "version"
	organizationIdColumnName  = "organization_id"
	creatorUsernameColumnName = "creator_username"
	tenderVersionIdColumnName = "tender_version_id"
	returningAllSuffix        = "RETURNING *"
	tenderAndVersionJoin      = versionTableName + " ON tender.tender_version_id = tender_version.id"
	selectTenderSum           = "tender.id, tender.status, tender_version.name, tender_version.description, " +
		"tender_version.service_type, tender_version.version, tender.organization_id, tender.creator_username, tender.created_at"
)

var (
	errTenderNotFound = fmt.Errorf("tender not found")
)

func NewTenderRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SaveTender(ctx context.Context, ten tender.Tender) (tender.Tender, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return tender.Tender{}, err
	}

	defer tx.Rollback(ctx)

	tenderBuilder := squirrel.Insert(tenderTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(statusColumnName, organizationIdColumnName, creatorUsernameColumnName).
		Values(ten.Status, ten.OrganizationId, ten.CreatorUsername).
		Suffix(returningAllSuffix)

	sql, args, err := tenderBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	savedTender, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Tender])
	if err != nil {
		fmt.Println("x1")
		return tender.Tender{}, err
	}

	versionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(tenderIdColumnName, serviceTypeColumnName, nameColumnName, descriptionColumnName, versionColumnName).
		Values(savedTender.Id, ten.ServiceType, ten.Name, ten.Description, 1).
		Suffix(returningAllSuffix)

	sql, args, err = versionBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	savedVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TenderVersion])
	if err != nil {
		return tender.Tender{}, err
	}

	setTenderVersion := squirrel.Update(tenderTableName).PlaceholderFormat(squirrel.Dollar).
		Set(tenderVersionIdColumnName, savedVersion.Id).
		Where(squirrel.Eq{idColumnName: savedTender.Id})

	sql, args, err = setTenderVersion.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	rows.Close()

	savedTender.TenderVersionId = sql2.NullInt32{
		Int32: int32(savedVersion.Id),
		Valid: true,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return tender.Tender{}, err
	}

	return model.DbTenderSumToTender(model.MergeTenderWithVersion(savedVersion, savedTender)), nil
}

func (r *repository) GetTenderById(ctx context.Context, id uuid.UUID) (tender.Tender, error) {
	op := "tender_repository.get_tender_by_id"
	builder := squirrel.Select(selectTenderSum).PlaceholderFormat(squirrel.Dollar).
		From(tenderTableName).Join(tenderAndVersionJoin).
		Where(squirrel.Eq{tenderTableName + "." + idColumnName: id})

	sql, args, err := builder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	fmt.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	version, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TenderSum])
	if err != nil {
		return tender.Tender{}, model2.NewNotFoundError(op, errTenderNotFound)
	}

	return model.DbTenderSumToTender(version), nil
}

func (r *repository) GetTenderList(ctx context.Context, page util.Page, serviceTypes []tender.ServiceType, username string, onlyPublished bool) ([]tender.Tender, error) {
	builder := squirrel.Select(selectTenderSum).PlaceholderFormat(squirrel.Dollar).
		From(tenderTableName).Join(tenderAndVersionJoin)

	if onlyPublished == true {
		builder = builder.Where(squirrel.Eq{statusColumnName: tender.Published})
	}

	if serviceTypes != nil && len(serviceTypes) > 0 {
		builder = builder.Where(squirrel.Eq{serviceTypeColumnName: serviceTypes})
	}

	if username != "" {
		builder = builder.Where(squirrel.Eq{creatorUsernameColumnName: username})
	}

	builder = builder.Offset(uint64(page.Offset)).Limit(uint64(page.Limit))

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	versions, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.TenderSum])
	if err != nil {
		return nil, err
	}

	return model.DdTenderVersionListToTenderList(versions), nil
}

func (r *repository) UpdateTenderStatus(ctx context.Context, id uuid.UUID, stat tender.Status) (tender.Tender, error) {
	updateBuilder := squirrel.Update(tenderTableName).PlaceholderFormat(squirrel.Dollar).Set(statusColumnName, stat).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		fmt.Println()
		return tender.Tender{}, err
	}

	fmt.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	rows.Close()

	return r.GetTenderById(ctx, id)
}

func (r *repository) UpdateTender(ctx context.Context, id uuid.UUID, name, description string, serviceType tender.ServiceType) (tender.Tender, error) {
	oldVersion, err := r.GetTenderById(ctx, id)
	if err != nil {
		return tender.Tender{}, err
	}

	setMap := make(map[string]interface{})

	setMap[versionColumnName] = oldVersion.Version + 1
	setMap[tenderIdColumnName] = oldVersion.Id.String()
	setMap[nameColumnName] = oldVersion.Name
	setMap[descriptionColumnName] = oldVersion.Description
	setMap[serviceTypeColumnName] = oldVersion.ServiceType

	if name != "" {
		setMap[nameColumnName] = name
	}

	if description != "" {
		setMap[descriptionColumnName] = description
	}

	if serviceType != "" {
		setMap[serviceTypeColumnName] = serviceType
	}

	newVersionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		SetMap(setMap).
		Suffix(returningAllSuffix)

	sql, args, err := newVersionBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	fmt.Println("sql:" + sql)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	newVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TenderVersion])
	if err != nil {
		return tender.Tender{}, err
	}

	updateTenderVersionIdBuilder := squirrel.Update(tenderTableName).PlaceholderFormat(squirrel.Dollar).
		Set(tenderVersionIdColumnName, newVersion.Id).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err = updateTenderVersionIdBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	fmt.Println("sql:" + sql)

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	rows.Close()

	oldVersion.Version = newVersion.Version
	oldVersion.Name = newVersion.Name
	oldVersion.Description = newVersion.Description
	oldVersion.ServiceType = tender.ServiceType(newVersion.ServiceType)

	return oldVersion, nil
}

func (r *repository) RollbackTender(ctx context.Context, id uuid.UUID, version int) (tender.Tender, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback(ctx)

	curTender, err := r.GetTenderById(ctx, id)
	if err != nil {
		return tender.Tender{}, err
	}

	getOldVersionBuilder := squirrel.Select("*").PlaceholderFormat(squirrel.Dollar).
		From(versionTableName).
		Where(squirrel.And{squirrel.Eq{tenderIdColumnName: id.String()}, squirrel.Eq{versionColumnName: version}})

	sql, args, err := getOldVersionBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	oldVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TenderVersion])
	if err != nil {
		return tender.Tender{}, err
	}

	versionBuilder := squirrel.Insert(versionTableName).PlaceholderFormat(squirrel.Dollar).
		Columns(tenderIdColumnName, serviceTypeColumnName, nameColumnName, descriptionColumnName, versionColumnName).
		Values(curTender.Id.String(), oldVersion.ServiceType, oldVersion.Name, oldVersion.Description, curTender.Version+1).
		Suffix(returningAllSuffix)

	sql, args, err = versionBuilder.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	savedVersion, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TenderVersion])
	if err != nil {
		return tender.Tender{}, err
	}

	setTenderVersion := squirrel.Update(tenderTableName).PlaceholderFormat(squirrel.Dollar).
		Set(tenderVersionIdColumnName, savedVersion.Id).
		Where(squirrel.Eq{idColumnName: id.String()})

	sql, args, err = setTenderVersion.ToSql()
	if err != nil {
		return tender.Tender{}, err
	}

	rows, err = r.pool.Query(ctx, sql, args...)
	if err != nil {
		return tender.Tender{}, err
	}

	rows.Close()

	curTender.ServiceType = tender.ServiceType(oldVersion.ServiceType)
	curTender.Name = oldVersion.Name
	curTender.Description = oldVersion.Description
	curTender.Version += 1

	err = tx.Commit(ctx)
	if err != nil {
		return tender.Tender{}, err
	}

	return curTender, nil
}
