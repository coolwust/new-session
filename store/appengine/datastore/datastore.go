package datastore

import (
	"time"
	"errors"
	"github.com/coldume/session"
	"github.com/coldume/session/store"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"

)

var ErrInvalidPropertyType = errors.New("datastore: invalid property type")

var _ store.Store = &Datastore{}

type Datastore struct {
	sessionKind, dataKind string
}

func NewDatastore(sessionKind, dataKind string) *Datastore {
	return &Datastore{sessionKind: sessionKind, dataKind: dataKind}
}

func (ds *Datastore) Get(ctx context.Context, id string) (*session.Session, error) {
	s := &Session{}
	k := datastore.NewKey(ctx, ds.sessionKind, id, 0, nil)
	if err := datastore.Get(ctx, k, s); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		}
		return nil, store.ErrNoSuchSession
	}
	d := make(Data)
	ite := datastore.NewQuery(ds.dataKind).Ancestor(k).Limit(1).Run(ctx)
	if _, err := ite.Next(d); err != nil {
		return nil, err
	}
	sess := session.NewSession(id, s.Duration, s.Expires)
	for k, v := range d {
		sess.Set(k, v)
	}
	return sess, nil
}

func (ds *Datastore) Put(ctx context.Context, sess *session.Session) error {
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		k := datastore.NewKey(ctx, ds.sessionKind, sess.ID, 0, nil)
		s := &Session{Duration: sess.Duration, Expires: sess.Expires}
		if _, err := datastore.Put(ctx, k, s); err != nil {
			return err
		}
		kk := datastore.NewIncompleteKey(ctx, ds.dataKind, k)
		d := make(Data)
		for n, v := range sess.All() {
			d[n] = v
		}
		if _, err := datastore.Put(ctx, kk, d); err != nil {
			return err
		}
		return nil
	}, nil)
	return err
}

func (ds *Datastore) Del(ctx context.Context, id string) error {
	return nil
}

func (ds *Datastore) Clear(ctx context.Context) error {
	return nil
}

func (ds *Datastore) Clean(ctx context.Context) error {
	return nil
}

type Session struct {
	Duration time.Duration `datastore:"Duration"`
	Expires  time.Time     `datastore:"Expires"`
}

var _ datastore.PropertyLoadSaver = make(Data)

type Data map[string]interface{}

func (data Data) Load(ps []datastore.Property) error {
	for _, p := range ps {
		data[p.Name] = p.Value
	}
	return nil
}

func (data Data) Save() ([]datastore.Property, error) {
	ps := make([]datastore.Property, 0)
	for n, v := range data {
		ps = append(ps, datastore.Property{Name: n, Value: v})
	}
	return ps, nil
}
