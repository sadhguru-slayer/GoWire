package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.RWMutex

	stop chan struct{}
	done chan struct{}
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}

	go func() {
		defer close(aof.done)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				aof.mu.Lock()
				aof.file.Sync()
				aof.mu.Unlock()

			case <-aof.stop:
				return
			}
		}
	}()

	return aof, nil
}

func (aof *Aof) CloseAof() error {
	close(aof.stop)
	<-aof.done

	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) WriteAof(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshall())
	if err != nil {
		return err
	}
	return nil
}

func (aof *Aof) ReadAof(callback func(value Value)) error {
	aof.mu.RLock()
	defer aof.mu.RUnlock()
	if _, err := aof.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	resp := NewRESP(aof.rd)

	for {
		value, err := resp.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		callback(value)
	}
	return nil
}
