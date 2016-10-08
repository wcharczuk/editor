package main

import "fmt"

// ANSI contains a bunch of ansi commands.
var ANSI = ansi{
	soh: byte(1),
	stx: byte(2),
	etx: byte(3),
	eot: byte(4),
	enq: byte(5),
	ack: byte(6),
	bs:  byte(8),
	tab: byte(9),
	lf:  byte(10),
	vt:  byte(11),
	cr:  byte(13),
	so:  byte(14),
	dle: byte(16),
	esc: byte(27),
	del: byte(127),

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
	soh byte
	stx byte
	etx byte
	eot byte
	enq byte
	ack byte
	bs  byte
	tab byte
	lf  byte
	vt  byte
	cr  byte
	esc byte
	so  byte
	dle byte
	del byte

	left  byte
	right byte
	up    byte
	down  byte

	escSequenceStart []byte
	clear            []byte
	clearLine        []byte
	hideCursor       []byte
	showCursor       []byte

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
