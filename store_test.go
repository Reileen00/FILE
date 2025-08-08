package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func newTestStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return newStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
func TestPathTransformFunc(t *testing.T) {
	key := "momomom"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "832ea9e4e91cb6dd11307273cd1f71d30a57d1e8"
	expectedPathName := "832ea/9e4e9/1cb6d/d1130/7273c/d1f71/d30a5/7d1e8"

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("Filename: have %s want %s", pathKey.Filename, expectedOriginalKey)
	}

	if pathKey.Pathname != expectedPathName {
		t.Errorf("Pathname: have %s want %s", pathKey.Pathname, expectedPathName)
	}

}

func TestStore(t *testing.T) {

	s := newTestStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d", i)

		data := []byte("some jpg bytes")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := ioutil.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("want %s have %s", string(data), string(b))
		}
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have that key %s", key)
		}
	}
}
