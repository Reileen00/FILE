package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrpytDecrypt(t *testing.T) {
	payload := "Foo not Bar"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}
	if nw != 11+len(payload) {
		t.Error(err)
	}
	fmt.Println(len(payload))
	fmt.Println(len(out.String()))
	fmt.Println(out.String())
	//fmt.Println(string(dst.Bytes()))
	fmt.Printf("raw bytes: %v\n", dst.Bytes())
	fmt.Printf("as string: %q\n", string(dst.Bytes()))
	fmt.Printf("raw bytes: %v\n", out.Bytes())
	fmt.Printf("as string: %q\n", string(out.Bytes()))

	if out.String() != payload {
		t.Errorf("decryption failed!!")
	}
}
