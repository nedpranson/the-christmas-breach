package user

import (
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/nedpranson/the-christmas-breach/core/db"
	"github.com/nedpranson/the-christmas-breach/core/util"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int64
	Username        string
	Seed            int64
	CompletedLevels uint16
	CreatedAt       time.Time
}

type CreateUserParams struct {
	Username string
	Pasword  string
}

const (
	MinUsernameLength = 1
	MaxUsernameLength = 32
)

func Create(ctx context.Context, p *CreateUserParams) (*User, error) {
	if err := ValidateUsername(p.Username); err != nil {
		return nil, err
	}

	seed := binary.LittleEndian.Uint64(util.CryptoRandomBytes(8))
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Pasword), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user, err := db.GetQueries().CreateUser(ctx, db.CreateUserParams{
		Username:     p.Username,
		PasswordHash: string(hash),
		Seed:         int64(seed),
	})

	if err, ok := errors.AsType[sqlite3.Error](err); ok {
		if err.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(err.Error(), "users.username") {
			return nil, ErrUsernameTaken
		}
	}

	if err != nil {
		return nil, err
	}

	return &User{
		ID:              user.ID,
		Username:        user.Username,
		Seed:            user.Seed,
		CompletedLevels: uint16(user.CompletedLevels),
		CreatedAt:       user.CreatedAt,
	}, nil
}

type Credentials struct {
	Username string
	Password string
}

func FromCredentials(ctx context.Context, creds *Credentials) (*User, error) {
	user, err := db.GetQueries().GetUserByUsername(ctx, creds.Username)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return &User{
		ID:              user.ID,
		Username:        user.Username,
		Seed:            user.Seed,
		CompletedLevels: uint16(user.CompletedLevels),
		CreatedAt:       user.CreatedAt,
	}, nil
}

func FromSession(ctx context.Context, session string) (*User, error) {
	user, err := db.GetQueries().GetUserFromSession(ctx, session)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidSession
	}

	if err != nil {
		return nil, err
	}

	return &User{
		ID:              user.ID,
		Username:        user.Username,
		Seed:            user.Seed,
		CompletedLevels: uint16(user.CompletedLevels),
		CreatedAt:       user.CreatedAt,
	}, nil
}

func (u *User) CreateSesion(ctx context.Context) (db.Session, error) {
	bytes := util.CryptoRandomBytes(32)
	token := hex.EncodeToString(bytes)

	session, err := db.GetQueries().CreateSession(ctx, db.CreateSessionParams{
		Token:     token,
		UserID:    u.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 14), // 2 weeks
	})

	if err != nil {
		return db.Session{}, err
	}

	return session, nil
}

func (u *User) IsLevelCompleted(level uint8) bool {
	return (u.CompletedLevels & 0x3FF & (1 << (level - 1))) == 1
}

func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return ErrInvalidUsername
	}

	for i := 0; i < len(username); i++ {
		if username[i] < '!' || username[i] > '~' {
			return ErrInvalidUsername
		}
	}

	return nil
}
