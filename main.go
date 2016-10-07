package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// Termios manages flags for the TTY.
type Termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]byte
	Ispeed uint64
	Ospeed uint64
}

const (
	getTermios = syscall.TIOCGETA
	setTermios = syscall.TIOCSETA
)

// TcSetAttr restores the terminal connected to the given file descriptor to a
// previous state.
func TcSetAttr(fd uintptr, termios *Termios) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(setTermios), uintptr(unsafe.Pointer(termios))); err != 0 {
		return err
	}
	return nil
}

// TcGetAttr retrieves the current terminal settings and returns it.
func TcGetAttr(fd uintptr) (*Termios, error) {
	var termios = &Termios{}
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, getTermios, uintptr(unsafe.Pointer(termios))); err != 0 {
		return nil, err
	}
	return termios, nil
}

// CfMakeRaw sets the flags stored in the termios structure to a state disabling
// all input and output processing, giving a ``raw I/O path''.
func CfMakeRaw(termios *Termios) {
	termios.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	termios.Oflag &^= syscall.OPOST
	termios.Lflag &^= (syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	termios.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0
}

// MakeRaw sets the flags stored in the termios structure for the given terminal fd
// to a state disabling all input and output processing, giving a ``raw I/O path''.
// It returns the current terminal's termios struct to allow to revert with TcSetAttr
func MakeRaw(fd uintptr) (*Termios, error) {
	old, err := TcGetAttr(fd)
	if err != nil {
		return nil, err
	}

	new := *old
	CfMakeRaw(&new)

	if err := TcSetAttr(fd, &new); err != nil {
		return nil, err
	}
	return old, nil
}

// ANSI contains a bunch of ansi commands.
var ANSI = ansi{
	endOfText: byte(3),

	start: byte(1),
	stx:   byte(2),
	end:   byte(5),
	ack:   byte(6),
	bs:    byte(8),
	tab:   byte(9),
	lf:    byte(10),
	vt:    byte(11),
	cr:    byte(13),
	so:    byte(14),
	dle:   byte(16),

	esc:           byte(27),
	del:           byte(127),
	sequenceStart: byte('['),

	escSequenceStart: []byte{byte(27), byte('[')},
	clear:            []byte{byte(27), byte('['), byte('2'), byte('J')},
	clearLine:        []byte{byte(27), byte('['), byte('2'), byte('K')},
	hideCursor:       []byte{byte(27), byte('['), byte('?'), byte('2'), byte('5'), byte('l')},
	showCursor:       []byte{byte(27), byte('['), byte('?'), byte('2'), byte('5'), byte('h')},

	colorReset:        []byte{byte(27), byte('['), byte('0'), byte('m')},
	colorBold:         []byte{byte(27), byte('['), byte('1'), byte('m')},
	colorItalics:      []byte{byte(27), byte('['), byte('1'), byte('m')},
	colorUnderline:    []byte{byte(27), byte('['), byte('1'), byte('m')},
	colorBoldOff:      []byte{byte(27), byte('['), byte('2'), byte('2'), byte('m')},
	colorItalicsOff:   []byte{byte(27), byte('['), byte('2'), byte('3'), byte('m')},
	colorUnderlineOff: []byte{byte(27), byte('['), byte('2'), byte('4'), byte('m')},
}

type ansi struct {
	endOfText byte
	start     byte
	stx       byte
	end       byte
	ack       byte
	bs        byte
	tab       byte
	lf        byte
	vt        byte
	cr        byte
	esc       byte
	so        byte
	dle       byte
	del       byte

	left  byte
	right byte
	up    byte
	down  byte

	escSequenceStart []byte
	clear            []byte
	clearLine        []byte
	hideCursor       []byte
	showCursor       []byte
	sequenceStart    byte

	colorReset        []byte
	colorBold         []byte
	colorBoldOff      []byte
	colorItalics      []byte
	colorItalicsOff   []byte
	colorUnderline    []byte
	colorUnderlineOff []byte
}

func (a ansi) Escape(sequence []byte) []byte {
	return append(a.escSequenceStart, sequence...)
}

func (a ansi) ClearScreen() []byte {
	return a.clear
}

func (a ansi) MoveCursor(row, col int) []byte {
	if row != 0 && col != 0 {
		return append(a.escSequenceStart, []byte(fmt.Sprintf("%d;%dH", row, col))...)
	} else if row != 0 {
		return append(a.escSequenceStart, []byte(fmt.Sprintf("%d;H", row))...)
	}
	return append(a.escSequenceStart, []byte(fmt.Sprintf(";%dH", col))...)
}

func (a ansi) Spaces(count int) []byte {
	switch count {
	case 0:
		return []byte{}
	case 1:
		return []byte{' '}
	case 2:
		return []byte{' ', ' '}
	case 3:
		return []byte{' ', ' ', ' '}
	case 4:
		return []byte{' ', ' ', ' ', ' '}
	case 5:
		return []byte{' ', ' ', ' ', ' ', ' '}
	case 6:
		return []byte{' ', ' ', ' ', ' ', ' ', ' '}
	case 7:
		return []byte{' ', ' ', ' ', ' ', ' ', ' '}
	case 8:
		return []byte{' ', ' ', ' ', ' ', ' ', ' ', ' '}
	case 9:
		return []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	case 10:
		return []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}
	var bytes []byte

	for x := 0; x < count; x++ {
		bytes = append(bytes, byte(' '))
	}
	return bytes
}

type cursor struct {
	row, col int
}

func (c cursor) Left() cursor {
	if c.col > 0 {
		return cursor{
			row: c.row,
			col: c.col - 1,
		}
	}
	return c
}

func (c cursor) LeftTab() cursor {
	newCol := c.col - 4
	if newCol < 0 {
		newCol = 0
	}
	return cursor{
		row: c.row,
		col: newCol,
	}
}

func (c cursor) Right() cursor {
	return cursor{
		row: c.row,
		col: c.col + 1,
	}
}

func (c cursor) RightTab() cursor {
	return cursor{
		row: c.row,
		col: c.col + 4,
	}
}

func (c cursor) Up() cursor {
	if c.row > 0 {
		return cursor{
			row: c.row - 1,
			col: c.col,
		}
	}
	return c
}

func (c cursor) Down() cursor {
	return cursor{
		row: c.row + 1,
		col: c.col,
	}
}

func (c cursor) BeginningOfLine() cursor {
	return cursor{
		row: c.row,
		col: 0,
	}
}

func (c cursor) DownBeginningOfLine() cursor {
	return cursor{
		row: c.row + 1,
		col: 0,
	}
}

type buffer [][]byte

func (b buffer) InsertRowAt(row int) buffer {
	if row >= len(b) {
		return append(b, []byte{})
	}

	output := make([][]byte, len(b)+1)
	afterInsert := false
	for y := 0; y < len(b); y++ {
		if y == row && !afterInsert {
			output[y] = []byte{}
			afterInsert = true
			y--
			continue
		}
		if afterInsert {
			output[y+1] = b[y][:]
		} else {
			output[y] = b[y][:]
		}
	}
	return output
}

func (b buffer) InsertCharacterAt(row, col int, c byte) buffer {
	if len(b) == 0 {
		return [][]byte{
			[]byte{c},
		}
	}
	output := make([][]byte, len(b))
	for y := 0; y < len(b); y++ {
		if y == row {
			if col >= len(b[y]) {
				output[y] = append(b[y], c)
			} else {
				output[y] = append(b[y][0:col], append([]byte{c}, b[y][col:]...)...) // zip
			}
		} else {
			output[y] = b[y][:]
		}
	}
	return output
}

func (b buffer) RemoveRowAt(row int) buffer {
	if len(b) == 0 {
		return b
	}
	output := make([][]byte, len(b)-1)
	afterSnip := false
	for y := 0; y < len(b); y++ {
		if y == row {
			afterSnip = true
			continue
		} else {
			if afterSnip {
				output[y-1] = b[y][:]
			} else {
				output[y] = b[y][:]
			}
		}
	}
	return output
}

func (b buffer) RemoveCharacterAt(row, col int) buffer {
	if len(b) == 0 {
		return b
	}

	if col < 0 {
		return b
	}

	output := make([][]byte, len(b))
	for y := 0; y < len(b); y++ {
		if y == row {
			if len(b[y]) > 0 {
				if col == len(b[y])-1 {
					output[y] = b[y][0:col]
				} else {
					output[y] = append(b[y][0:col], b[y][col+1:]...) // snip
				}
			}
		} else {
			output[y] = b[y][:]
		}
	}
	return output
}

func (b buffer) TrimRowAt(row, col int) buffer {
	if len(b) == 0 {
		return b
	}

	output := make([][]byte, len(b))
	for y := 0; y < len(b); y++ {
		if y == row {
			if len(b[y]) > 0 {
				output[y] = b[y][0:col]
			}
		} else {
			output[y] = b[y][:]
		}
	}
	return output
}

type editorState struct {
	buffer buffer
	scroll int //denotes scrollTop, or where we start drawing the buffer
	cursor cursor
}

func (es editorState) Write(b byte) editorState {
	return editorState{
		buffer: es.buffer.InsertCharacterAt(es.cursor.row, es.cursor.col, b),
		cursor: es.cursor.Right(),
	}
}

func (es editorState) MoveLeft() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.Left(),
	}
}

func (es editorState) MoveRight() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.Right(),
	}
}

func (es editorState) MoveUp() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.Up(),
	}
}

func (es editorState) MoveDown() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.Down(),
	}
}

func (es editorState) MoveToBeginningOfLine() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.BeginningOfLine(),
	}
}

func (es editorState) MoveToEndOfLine() editorState {
	return editorState{
		buffer: es.buffer,
		cursor: cursor{
			row: es.cursor.row,
			col: len(es.buffer[es.cursor.row]),
		},
	}
}

func (es editorState) Newline() editorState {
	//newline does the following, it both creates a new line, but also does a bunch of line manipulation
	// - on an existing line, if there is text after the cursor
	//		which pushes existing content down one line
	// - on an empty line i just creates a new line

	if len(es.buffer[es.cursor.row]) == 0 {
		return editorState{
			buffer: es.buffer.InsertRowAt(es.cursor.row + 1),
			cursor: es.cursor.DownBeginningOfLine(),
		}
	}

	if es.cursor.col == len(es.buffer[es.cursor.row]) { // if we're at the end of the row
		return editorState{
			buffer: es.buffer.InsertRowAt(es.cursor.row + 1),
			cursor: es.cursor.DownBeginningOfLine(),
		}
	}
	return es
}

func (es editorState) Backspace() editorState {
	// nuke the character at the cursor, move the cursor to the left

	// if we're at the beginning of a line
	if es.cursor.col == 0 {

		// and we're on the first line
		if es.cursor.row == 0 {
			return es //just return the state
		}

		// else move up a row, to the end of the line
		newRow := es.cursor.row - 1
		return editorState{
			buffer: es.buffer,
			cursor: cursor{
				row: newRow,
				col: len(es.buffer[newRow]),
			},
		}
	}

	// remove the previous character per normal.
	return editorState{
		buffer: es.buffer.RemoveCharacterAt(es.cursor.row, es.cursor.col-1),
		cursor: es.cursor.Left(),
	}
}

func (es editorState) TrimLine() editorState {
	if es.cursor.col == 0 {
		if es.cursor.row == 0 {
			return editorState{
				buffer: es.buffer.RemoveRowAt(es.cursor.row),
				cursor: es.cursor,
			}
		}
		return editorState{
			buffer: es.buffer.RemoveRowAt(es.cursor.row),
			cursor: cursor{
				row: es.cursor.row - 1,
				col: len(es.buffer[es.cursor.row-1]),
			},
		}
	}
	return editorState{
		buffer: es.buffer.TrimRowAt(es.cursor.row, es.cursor.col),
		cursor: es.cursor,
	}
}

func processInput(b byte, state editorState) (editorState, error) {
	switch b {
	case ANSI.endOfText:
		return state, errors.New("should exit")
	case ANSI.vt:
		return state.TrimLine(), nil
	case ANSI.dle:
		return state.MoveUp(), nil
	case ANSI.so:
		return state.MoveDown(), nil
	case ANSI.start:
		return state.MoveToBeginningOfLine(), nil
	case ANSI.stx:
		return state.MoveLeft(), nil
	case ANSI.ack:
		return state.MoveRight(), nil
	case ANSI.end:
		return state.MoveToEndOfLine(), nil
	case ANSI.bs, ANSI.del:
		return state.Backspace(), nil
	case ANSI.cr, ANSI.lf:
		return state.Newline(), nil
	default:
		return state.Write(b), nil
	}
}

func render(tty *os.File, state editorState) (err error) {
	tty.Write(ANSI.ClearScreen())
	tty.Write(ANSI.MoveCursor(0, 0))
	tty.Write(ANSI.colorReset)

	c := make([]byte, 1)
	var cursorRowTabs int
	for row := 0; row < len(state.buffer); row++ {
		tty.Write(ANSI.MoveCursor(row+1, 0))
		for col := 0; col < len(state.buffer[row]); col++ {
			c[0] = state.buffer[row][col]
			switch c[0] {
			case ANSI.tab:
				tty.Write(ANSI.Spaces(4))
				if row == state.cursor.row {
					cursorRowTabs++
				}
			default:
				tty.Write(c)
			}
		}
		if row != state.cursor.row {
			cursorRowTabs = 0
		}
	}

	var extraTabSpaces int
	if cursorRowTabs > 0 {
		extraTabSpaces = (cursorRowTabs * 3)
	}

	tty.Write(ANSI.MoveCursor(state.cursor.row+1, state.cursor.col+1+extraTabSpaces))
	return
}

func initTerm() (*Termios, *os.File) {
	tty := os.Stdout
	initialSettings, err := MakeRaw(os.Stdout.Fd())
	if err != nil {
		log.Fatal(err)
	}

	tty.Write(ANSI.ClearScreen())
	tty.Write(ANSI.MoveCursor(0, 0))

	return initialSettings, tty
}

func restoreTerm(initialSettings *Termios, tty *os.File) {
	err := TcSetAttr(tty.Fd(), initialSettings)
	tty.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	initialSettings, tty := initTerm()
	defer restoreTerm(initialSettings, tty)

	state := editorState{
		buffer: [][]byte{
			[]byte{},
		},
	}
	var c = make([]byte, 1)
	for {
		os.Stdin.Read(c)
		state, err = processInput(c[0], state)
		if err != nil {
			return
		}
		render(tty, state)
	}
}
