package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/nurkenti/hadiaParfums/internal/repository/sqlc"
)

type AdminService struct {
	store *db.Queries
}

func NewAdminService(store *db.Queries) *AdminService {
	return &AdminService{store: store}
}

// Создаем админа
func (a *AdminService) CreateAdmin(ctx context.Context, chatID int64, name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	_, err := a.store.CreateUser(ctx, db.CreateUserParams{
		ChatID:  chatID,
		Lang:    pgtype.Text{String: "ru", Valid: true},
		Name:    name,
		IsAdmin: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return errors.New("fuck adm")
	}

	return nil
}

func (a *AdminService) MakeAdmin(ctx context.Context, chatID int64) error {
	_, err := a.store.UpdateUser(ctx, db.UpdateUserParams{
		ChatID: chatID,
		Lang:   pgtype.Text{String: "ru", Valid: true},
	})
	return err
}

// Проверям на админ.
func (a *AdminService) IsAdmin(ctx context.Context, chatID int64) (bool, error) {
	user, err := a.store.GetUser(ctx, chatID)
	if err != nil {
		return false, nil
	}
	if user.IsAdmin.Valid {
		return user.IsAdmin.Bool, nil
	}
	return false, nil
}
