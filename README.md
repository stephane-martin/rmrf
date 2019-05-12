# rmrf

Sometimes "rm -rf" in a shell will fail/hang/takes forever when the file tree
you try to remove has millions of entries.

rmrf is just a Go implementation of "rm -rf" that tries to be fast and efficient.

## Usage

- `rmrf`: delete everything in current directory. Basically "rm -rf *"
- `rmrf PATH [PATH2...]`: delete the provided paths

The implementation does not use recursive functions to walk the file tree.

## Install

go get -u github.com/stephane-martin/rmrf

