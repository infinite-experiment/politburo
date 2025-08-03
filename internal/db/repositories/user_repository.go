package repositories

import (
	"context"
	"log"

	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/entities"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *entities.User) error {

	return r.db.QueryRowxContext(ctx, constants.InsertUser,
		user.DiscordID,
		user.IFCommunityID,
		user.IFApiID,
		user.IsActive,
	).StructScan(user)
}

func (r *UserRepository) FindUserByDiscordId(ctx context.Context, discordId string) (*entities.User, error) {

	var user entities.User

	err := r.db.QueryRowxContext(ctx, constants.GetUserByDiscordId, discordId).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteAllUsers(ctx context.Context) error {
	err2 := r.db.QueryRowxContext(ctx, constants.DeleteAllRoles)
	err1 := r.db.QueryRowxContext(ctx, constants.DeleteAllServers)
	err3 := r.db.QueryRowxContext(ctx, constants.DeleteAllUsers)
	log.Printf("Query output: %v \n ===== \n %v \n =========\n %v", err1, err2, err3)
	return nil
}

func (r *UserRepository) FindUserMembership(ctx context.Context, sDiscordId string, uDiscordId string) (*entities.Membership, error) {
	var membership entities.Membership

	if err := r.db.GetContext(ctx, &membership, constants.GetUserMembership, sDiscordId, uDiscordId); err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *UserRepository) InsertMembership(
	ctx context.Context,
	m *entities.UserVARole,
) error {
	return r.db.QueryRowxContext(
		ctx,
		constants.InsertMembership,
		m.UserID,
		m.VAID,
		m.Role, // enums.RoleAdmin etc.
	).Scan(&m.ID, &m.JoinedAt)
}
