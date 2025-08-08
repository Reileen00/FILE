package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

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

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momoijdfif"

	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momoijdfif"

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
	s.Delete(key)
}
