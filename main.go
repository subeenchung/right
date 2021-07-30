package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"sync"
)

func main() {
	fname := os.Args[1]
	h, err := getFileHash(fname)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	fmt.Printf("SHA256 hash of file: %s\n%s\n", fname, base64.StdEncoding.EncodeToString(h))
}

func getFileHash(path string) ([]byte, error) {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	wg.Add(1)
	h := sha256.New()

	go func(h hash.Hash) hash.Hash {
		fh, err := os.Open(path)
		if err != nil {
			errCh <- err
		}
		if _, err := io.Copy(h, fh); err != nil {
			errCh <- err
		}
		wg.Done()
		return h
	}(h)
	wg.Wait()
	if len(errCh) != 0 {
		return []byte{}, <-errCh
	}
	return h.Sum(nil), nil
}
