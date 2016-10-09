package main

import (
	"fmt"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestBufferInsertRowAt(t *testing.T) {
	assert := assert.New(t)

	b := buffer{}
	edited := b.InsertRowAt(0)
	assert.Len(edited, 1)
}

func TestBufferInsertRowAtMid(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a'},
		[]byte{'b'},
	}
	edited := b.InsertRowAt(1)
	assert.Len(edited, 3)
	assert.Equal('a', edited[0][0])
	assert.Len(edited[1], 0)
	assert.Len(edited[2], 1)
	assert.Equal('b', edited[2][0])
}

func TestBufferInsertRowAtEnd(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a'},
		[]byte{'b'},
	}
	edited := b.InsertRowAt(2)
	assert.Len(edited, 3)
	assert.Equal('a', edited[0][0])
	assert.Len(edited[1], 1)
	assert.Equal('b', edited[1][0])
	assert.Len(edited[2], 0)

	edited = edited.InsertRowAt(3)
	assert.Len(edited, 4)
	assert.Equal('a', edited[0][0])
	assert.Len(edited[1], 1)
	assert.Equal('b', edited[1][0])
	assert.Len(edited[2], 0)
	assert.Len(edited[3], 0)
}

func TestBufferMoveAfterToNextRow(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c', 'd'},
		[]byte{'z', 'x', 'y'},
	}

	edited := b.MoveAfterToNextRow(0, 2)
	assert.Len(edited, 3)
	assert.Len(edited[0], 2)
	assert.Len(edited[1], 2)
	assert.Len(edited[2], 3)

	assert.Equal('b', edited[0][1])
	assert.Equal('c', edited[1][0])
	assert.Equal('z', edited[2][0])
}

func TestBufferMoveRowToEndOfPrevious(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c', 'd'},
		[]byte{'z', 'x', 'y'},
	}

	edited := b.MoveRowToEndOfPrevious(1)
	assert.Len(edited, 1)
	assert.Len(edited[0], 7)

	assert.Equal('a', edited[0][0])
	assert.Equal('y', edited[0][6])
}

func TestInsertAtEmptyLine(t *testing.T) {
	assert := assert.New(t)

	b := buffer{}
	edited := b.InsertCharacterAt(0, 0, 'a')
	assert.Len(edited, 1)
	assert.Len(edited[0], 1)
}

func TestInsertAtEndOfLine(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a'},
	}
	edited := b.InsertCharacterAt(0, 1, 'b')
	assert.Len(edited, 1)
	assert.Len(edited[0], 2, fmt.Sprintf("%#v", edited[0]))
	assert.Equal('b', edited[0][1])
}

func TestInsertMidLine(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c'},
	}
	edited := b.InsertCharacterAt(0, 1, 'd')
	assert.Len(edited, 1)
	assert.Len(edited[0], 4)
	assert.Equal('a', edited[0][0])
	assert.Equal('d', edited[0][1])
	assert.Equal('b', edited[0][2])
}

func TestBufferRemoveCharacterAtStart(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c'},
		[]byte{'a', 'b', 'c'},
		[]byte{'a', 'b', 'c'},
	}
	edited := b.RemoveCharacterAt(1, 0)
	assert.Len(edited, 3)
	assert.Len(edited[0], 3)
	assert.Len(edited[1], 2)
	assert.Len(edited[2], 3)
	assert.Equal('b', edited[1][0])
	assert.Equal('c', edited[1][1])
}

func TestBufferRemoveCharacterAtMid(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c'},
	}
	edited := b.RemoveCharacterAt(0, 1)
	assert.Len(edited, 1)
	assert.Len(edited[0], 2)
	assert.Equal('c', edited[0][1])
}

func TestBufferRemoveCharacterAtEnd(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c'},
	}
	edited := b.RemoveCharacterAt(0, 2)
	assert.Len(edited, 1)
	assert.Len(edited[0], 2)
	assert.Equal('b', edited[0][1])
}

func TestBufferTrimLineAt(t *testing.T) {
	assert := assert.New(t)

	var b buffer = [][]byte{
		[]byte{'a', 'b', 'c', 'd', 'e'},
		[]byte{'a', 'b', 'c', 'd', 'e'},
		[]byte{'a', 'b', 'c', 'd', 'e'},
	}

	edited := b.TrimRowAt(1, 2)
	assert.Len(edited, 3)
	assert.Len(edited[0], 5)
	assert.Len(edited[1], 2)
	assert.Len(edited[2], 5)
}
