package main

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
	if es.cursor.col == len(es.buffer[es.cursor.row]) {
		return es
	}
	return editorState{
		buffer: es.buffer,
		cursor: es.cursor.Right(),
	}
}

func (es editorState) MoveUp() editorState {
	if es.cursor.row == 0 {
		return es
	}

	previousLineLength := len(es.buffer[es.cursor.row-1])
	if es.cursor.col < previousLineLength {
		return editorState{
			buffer: es.buffer,
			cursor: es.cursor.Up(),
		}
	}
	return editorState{
		buffer: es.buffer,
		cursor: cursor{
			row: es.cursor.row - 1,
			col: previousLineLength,
		},
	}
}

func (es editorState) MoveDown() editorState {
	if es.cursor.row == len(es.buffer)-1 {
		return es
	}

	nextLineLength := len(es.buffer[es.cursor.row+1])

	if es.cursor.col < nextLineLength {
		return editorState{
			buffer: es.buffer,
			cursor: es.cursor.Down(),
		}
	}
	return editorState{
		buffer: es.buffer,
		cursor: cursor{
			row: es.cursor.row + 1,
			col: nextLineLength,
		},
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

	return editorState{
		buffer: es.buffer.MoveAfterToNextRow(es.cursor.row, es.cursor.col),
		cursor: es.cursor.DownBeginningOfLine(),
	}
}

func (es editorState) Backspace() editorState {

	if es.cursor.col == 0 { // if we're at the beginning of a line
		if es.cursor.row == 0 { // and we're on the first line
			return es //just return the state
		}

		// else move up a row, to the end of the line
		previousRow := es.cursor.row - 1
		endOfLine := len(es.buffer[previousRow])

		if len(es.buffer[es.cursor.row]) == 0 {
			return editorState{
				buffer: es.buffer.RemoveRowAt(es.cursor.row),
				cursor: cursor{
					row: previousRow,
					col: endOfLine,
				},
			}
		}

		// move up a line.
		return editorState{
			buffer: es.buffer.MoveRowToEndOfPrevious(es.cursor.row),
			cursor: cursor{
				row: previousRow,
				col: endOfLine,
			},
		}
	}

	// nuke the character at the cursor, move the cursor to the left
	return editorState{
		buffer: es.buffer.RemoveCharacterAt(es.cursor.row, es.cursor.col-1),
		cursor: es.cursor.Left(),
	}
}

func (es editorState) TrimLine() editorState {
	return editorState{
		buffer: es.buffer.TrimRowAt(es.cursor.row, es.cursor.col),
		cursor: es.cursor,
	}
}
