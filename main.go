package main

import (
	"errors"
	"log"
	"os"
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
		if c[0] == 0x1b { // special key
			os.Stdin.Read(c)
			if c[0] == 0x5b { // escape
				os.Stdin.Read(c)
				switch c[0] { // arrow direction
				case 0x43:
					state = state.MoveRight()
					render(tty, state)
					continue
				case 0x44:
					state = state.MoveLeft()
					render(tty, state)
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
			render(tty, state)
		}
	}
}
