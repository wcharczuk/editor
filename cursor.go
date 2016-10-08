package main

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
