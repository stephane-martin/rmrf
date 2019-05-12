package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type entry struct {
	path   string
	parent *entry
	file   *os.File
}

var pool *sync.Pool

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return new(entry)
		},
	}
}

func newEntry(path string, parent *entry, f *os.File) *entry {
	e := pool.Get().(*entry)
	e.path = path
	e.file = f
	e.parent = parent
	return e
}

func releaseEntry(e *entry) {
	pool.Put(e)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		d, err := os.Open(".")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer d.Close()
		for {
			names, err := d.Readdirnames(1)
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			rmrf(names[0])
		}
	}
	for _, arg := range args {
		rmrf(arg)
	}
}

func (e *entry) next() *entry {
	for {
		infos, err := e.file.Readdir(1)
		if err == io.EOF {
			_ = e.file.Close()
			return nil
		}
		if err != nil {
			_ = e.file.Close()
			fmt.Fprintf(os.Stderr, "readdir %s: %s\n", e.path, err)
			return nil
		}
		fullpath := filepath.Join(e.path, infos[0].Name())

		if !infos[0].IsDir() {
			err := os.Remove(fullpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "remove file %s: %s\n", fullpath, err)
			}
			continue
		}

		f, err := os.Open(fullpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %s\n", fullpath, err)
			continue
		}
		return newEntry(fullpath, e, f)
	}
}

func rmrf(name string) {
	info, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stat %s: %s\n", name, err)
		return
	}
	if !info.IsDir() {
		err := os.Remove(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "remove file %s: %s\n", name, err)
		}
		return
	}
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open %s: %s\n", name, err)
		return
	}
	current := newEntry(name, nil, f)
	for {
		e := current.next()
		if e != nil {
			current = e
			continue
		}
		err := os.Remove(current.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "remove directory %s: %s\n", current.path, err)
		}
		parent := current.parent
		if parent == nil {
			return
		}
		releaseEntry(current)
		current = parent
	}
}
