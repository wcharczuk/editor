package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
)

const (
	byteNewLine = byte('\n')
	byteTab     = byte('\t')
)

func processSingleInput(b byte, state editorState) (editorState, error) {
	switch b {
	case ANSI.etx:
		return state, errors.New("should exit")
	case ANSI.vt:
		return state.TrimLine(), nil
	case ANSI.dle:
		return state.MoveUp(), nil
	case ANSI.so:
		return state.MoveDown(), nil
	case ANSI.soh:
		return state.MoveToBeginningOfLine(), nil
	case ANSI.stx:
		return state.MoveLeft(), nil
	case ANSI.ack:
		return state.MoveRight(), nil
	case ANSI.enq:
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

func stateFromFile(path string) (editorState, error) {
	f, err := os.Open(path)
	if err != nil {
		return editorState{}, err
	}
	return stateFromReader(f), nil
}

func stateFromReader(reader io.ReaderAt) editorState {
	es := editorState{
		buffer: [][]byte{},
	}

	var cursor int64
	var readBuffer = make([]byte, 32)
	var readErr error
	var lineBuffer = bytes.NewBuffer([]byte{})
	for readErr == nil {
		lineBuffer.Reset()
		cursor, readErr = readLine(reader, cursor, readBuffer, lineBuffer)
		if readErr != nil {
			continue
		}
		es.buffer = append(es.buffer, lineBuffer.Bytes())
	}
	return es
}

// readLine reads a file until a newline.
func readLine(f io.ReaderAt, cursor int64, readBuffer []byte, lineBuffer *bytes.Buffer) (int64, error) {
	// bytesRead is the return from the ReadAt function
	// it indicates how many effective bytes we read from the stream.
	var bytesRead int
	// err is our primary indicator if there was an issue with the stream
	// or if we've reached the end of the file.
	var err error
	// b is the byte we're reading at a time.
	var b byte

	// while we haven't hit an error (this includes EOF!)
	for err == nil {
		// read the stream
		bytesRead, err = f.ReadAt(readBuffer, cursor)
		// abort on error
		if err != nil && err != io.EOF { //let this continue on eof
			return cursor, err
		}

		// slurp the read buffer.
		for readBufferIndex := 0; readBufferIndex < bytesRead; readBufferIndex++ {
			// advance the cursor regardless of what we read out.
			// if we read a newline, great! we'll start the next character after the newline after.
			cursor++

			// slurp the byte out of the read buffer
			b = readBuffer[readBufferIndex]
			if b == byteNewLine {
				// we bifurcate here because we need to forward the eof
				// if we read the buffer exactly right.
				if readBufferIndex == bytesRead-1 {
					return cursor, err
				}
				// otherwise the newline may have happened
				// before the actual eof.
				return cursor, nil
			}

			// b wasnt a newline, write it to the output buffer.
			lineBuffer.WriteByte(b)
		}
	}
	// we've reached the end of the file
	// there may not have been a newline
	// return what we have
	return cursor, err
}

func main() {
	var err error

	var state editorState
	if len(os.Args) < 2 {
		state = newEditorState()
	} else {
		state, err = stateFromFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	initialSettings, tty := initTerm()
	defer restoreTerm(initialSettings, tty)

	var c = make([]byte, 1)
	for {
		render(tty, state)

		os.Stdin.Read(c)
		if c[0] == 0x1b { // special key
			os.Stdin.Read(c)
			if c[0] == 0x5b { // escape
				os.Stdin.Read(c)
				switch c[0] { // arrow direction
				case 0x43:
					state = state.MoveRight()
					continue
				case 0x44:
					state = state.MoveLeft()
					continue
				default:
					continue
				}
			}
		} else { //normal input
			state, err = processSingleInput(c[0], state)
			if err != nil {
				return
			}
		}
	}
}
