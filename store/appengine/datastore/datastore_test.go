package datastore

import (
	"testing"
	"reflect"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/aetest"
)

var dataTest = Data{
	"foo": int64(1),
	"bar": float64(1.5),
	"baz": "qux",
}

func TestData(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	k := datastore.NewIncompleteKey(ctx, "data", nil)
	if k, err = datastore.Put(ctx, k, dataTest); err != nil {
		t.Fatal(err)
	}
	d := make(Data)
	if err := datastore.Get(ctx, k, d); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dataTest, d) {
		t.Fatalf("got %v, want %v", d, dataTest)
	}
}
