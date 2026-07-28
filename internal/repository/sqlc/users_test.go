package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

var id = 324

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		ChatID: 12315233442,
		Lang:   pgtype.Text{String: "kz"},
		Name:   "Saule",
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		require.NoError(t, err)
		require.NotEmpty(t, user)

		require.Equal(t, arg.ChatID, user.ChatID)
		require.Equal(t, arg.Lang.String, user.Lang.String)

		require.NotZero(t, arg.ChatID, user.ChatID)
		require.NotZero(t, arg.Lang.String, user.Lang.String)
		require.Empty(t, arg.Lang.String)
	}
} 