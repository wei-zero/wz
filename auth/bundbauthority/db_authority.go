package bundbauthority

import (
	"context"
	"errors"

	"math/rand"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/wei-zero/wz/auth"
)

const (
	tokenStrLen = 20
)

type authority struct {
	db       *bun.DB
	expireIn int64
}

type Token struct {
	bun.BaseModel `bun:"table:tokens"`

	AccessToken string `bun:"type:varchar(60),pk"`
	UserId      int64
	Roles       string
	ExpiredAt   time.Time
}

func NewAuthority(db *bun.DB, expireInSec int64) auth.Authority {
	return &authority{
		db:       db,
		expireIn: expireInSec,
	}
}

func (a *authority) Revoke(tokenStr string) error {
	token := Token{AccessToken: tokenStr}
	_, err := a.db.NewDelete().Model(&token).WherePK().Exec(context.Background())
	return err
}

func (a *authority) Issue(identity *auth.Identity, duration time.Duration) (string, int64, error) {
	if duration == 0 {
		duration = time.Duration(a.expireIn) * time.Second
	}

	expiresAt := time.Now().Add(duration)
	// generate token
	accessToken := generateRandomStr(tokenStrLen)

	_, err := a.db.NewInsert().Model(&Token{
		AccessToken: accessToken,
		UserId:      identity.UserID,
		Roles:       strings.Join(identity.Roles, ","),
		ExpiredAt:   expiresAt,
	}).Exec(context.Background())

	if err != nil {
		return "", 0, err
	}

	return accessToken, expiresAt.UnixNano() / 1e6, nil
}

func generateRandomStr(tokenStrLen int) string {
	randStr := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(randStr)
	result := []byte{}
	for i := 0; i < tokenStrLen; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result)
}

func (a *authority) Verify(tokenStr string) (*auth.Identity, int64, error) {
	if len(tokenStr) != tokenStrLen {
		return nil, 0, errors.New("bad token")
	}
	var t Token
	err := a.db.NewSelect().Model(&t).Where("access_token", tokenStr).Scan(context.Background())
	if err != nil {
		return nil, 0, err
	}

	if time.Now().After(t.ExpiredAt) {
		return nil, 0, auth.ErrTokenExpired
	}

	return &auth.Identity{
		UserID: t.UserId,
		Roles:  strings.Split(t.Roles, ","),
	}, t.ExpiredAt.UnixNano() / 1e6, nil
}
