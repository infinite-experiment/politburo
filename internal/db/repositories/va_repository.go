package repositories

import (
	"context"
	"fmt"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/entities"

	"github.com/jmoiron/sqlx"
)

type VARepository struct {
	db *sqlx.DB
}

func NewVARepository(db *sqlx.DB) *VARepository {
	return &VARepository{db}
}

func (r *VARepository) InsertVA(ctx context.Context, va *entities.VA) error {
	return r.db.QueryRowxContext(
		ctx,
		constants.InsertVA,
		va.Name,
		va.Code,
		va.DiscordID,
		va.CallsignPrefix,
		va.CallsignSuffix,
		va.IsActive).StructScan(va)
}

func (r *VARepository) FindVAByDiscordServerID(ctx context.Context, sID string) (*entities.VA, error) {
	var va entities.VA
	err := r.db.QueryRowxContext(ctx, constants.GetVAByDiscordID, sID).StructScan(&va)
	if err != nil {
		return nil, err
	}

	return &va, nil
}

func (r *VARepository) InsertVAWithAdmin(
	ctx context.Context,
	va *entities.VA, // ID & timestamps filled on success
	adminID string, // UUID of the requesting user
) (*entities.UserVARole, error) {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // safe even after Commit

	/* --- A. insert virtual_airlines row -------------------------- */
	if err := sqlx.GetContext(ctx, tx, va, constants.InsertVA,
		va.Name,
		va.Code,
		va.DiscordID,
		va.CallsignPrefix,
		va.CallsignSuffix,
		va.IsActive,
	); err != nil {
		return nil, fmt.Errorf("insert VA: %w", err)
	}

	/* --- B. insert admin membership row -------------------------- */
	m := &entities.UserVARole{
		UserID:   adminID,
		VAID:     va.ID,
		Role:     constants.RoleAdmin,
		IsActive: true,
	}

	if err := tx.QueryRowxContext(
		ctx,
		constants.InsertMembership, // already defined in constants
		m.UserID,
		m.VAID,
		m.Role,
	).Scan(&m.ID, &m.JoinedAt); err != nil {
		return nil, fmt.Errorf("insert membership: %w", err)
	}

	/* --- Commit -------------------------------------------------- */
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return m, nil
}
