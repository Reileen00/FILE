package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "ggnetwork"

// CASPathTransformFunc converts a key to a hashed path.
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath(root string) string {
	return root + "/" + p.Pathname + "/" + p.Filename
}

// FirstPathName returns only the first folder component of the path.
func (p PathKey) FirstPathName() string {
	parts := strings.Split(p.Pathname, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func newStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := pathKey.FullPath(s.Root) // ✅ fixed: correctly pass root
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

// FIXME: Instead of copying directly to a reader we first copy
// this into a buffre. MAybe just return the File from readStream ?
func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	file, err := os.Open(pathKey.FullPath(s.Root))
	if err != nil {
		return 0, nil, err
	}
	fi, err := file.Stat()

	if err != nil {
		return 0, nil, err
	}
	return fi.Size(), file, nil
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)

	// Ensure the directory exists
	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return 0, err
	}

	fullPathWithRoot := pathKey.FullPath(s.Root)
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}
	return n, nil
}
