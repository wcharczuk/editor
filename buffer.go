package main

type buffer [][]byte

func (b buffer) RowLength(row int) int {
	if len(b) == 0 {
		return 0
	}
	return len(b[row])
}

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

func (b buffer) MoveAfterToNewRow(row, col int) buffer {
	if row >= len(b) {
		return append(b, []byte{})
	}

	output := make([][]byte, len(b)+1)
	for y := 0; y <= len(b); y++ { // <= means extra row
		if y == row {
			output[y] = b[y][0:col]
			continue
		}
		if y == row+1 {
			output[y] = b[y-1][col:]
			continue
		}
		if y > row {
			output[y] = b[y-1][:]
		} else {
			output[y] = b[y][:]
		}
	}
	return output
}
