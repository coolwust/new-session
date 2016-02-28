package store

import (
	"errors"
	"github.com/coldume/session"
	"golang.org/x/net/context"
)

type Store interface {

	// should return ErrNoSuchSession when no session is found
	Get(ctx context.Context, id string) (*session.Session, error)

	Put(ctx context.Context, sess *session.Session) error

	Del(ctx context.Context, id string) error

	Clear(ctx context.Context) error

	Clean(ctx context.Context) error
}

var ErrNoSuchSession = errors.New("session: no such session")
