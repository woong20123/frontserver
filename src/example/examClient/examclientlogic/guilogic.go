package examclientlogic

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var editBox EditBox
var readHigh int

const editBoxWidth = 100

var readStrings []readString

type readString struct {
	msg     string
	wordcol termbox.Attribute
}

func runReadLogic() {
	readEvent := Instance().ObjMgr().ChanManager().ChanRequestToGui()
	for {
		select {
		case e := <-readEvent:
			switch e.Type {
			case ToGUIEnum.TYPEMsgPrint:
				insertReadmsg(e.Msg, e.wordcol)
			case ToGUIEnum.TYPEWindowClear:
				clearReadmsg()
			}

			redrawAll()
		}
	}
}

func insertReadmsg(msg string, wordcol termbox.Attribute) {
	readStrings = append(readStrings, readString{msg, wordcol})
	strlength := len(readStrings)
	readSize := readHigh - 1
	if strlength >= readSize {
		readStrings = readStrings[strlength-readSize:]
	}
}

func clearReadmsg() {
	strlength := len(readStrings)
	readStrings = readStrings[strlength:]
}

func printReadmsg() {
	const coldef = termbox.ColorDefault
	w, _ := termbox.Size()
	xpos := ((w / 2) - (editBoxWidth / 2))
	ypos := 1
	for index, value := range readStrings {
		tbprint(xpos, ypos+index, value.wordcol, coldef, value.msg)
	}
}

func drawBoxUI(xpos int, ypos int, width int, high int) {
	const coldef = termbox.ColorDefault
	widthhalf := width / 2
	highhalf := high / 2
	fill(xpos-widthhalf, ypos-highhalf, 1, high, termbox.Cell{Ch: '|'})
	fill(xpos+widthhalf, ypos-highhalf, 1, high, termbox.Cell{Ch: '|'})
	fill(xpos-widthhalf, ypos-highhalf, width, 1, termbox.Cell{Ch: '-'})
	fill(xpos-widthhalf, ypos+highhalf, width, 1, termbox.Cell{Ch: '-'})
	termbox.SetCell(xpos-widthhalf, ypos-highhalf, '+', coldef, coldef)
	termbox.SetCell(xpos+widthhalf, ypos-highhalf, '+', coldef, coldef)
	termbox.SetCell(xpos-widthhalf, ypos+highhalf, '+', coldef, coldef)
	termbox.SetCell(xpos+widthhalf, ypos+highhalf, '+', coldef, coldef)
}

func redrawAll() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	readHigh = (h * 80) / 100
	readBoxmidx := w / 2
	readBoxmidy := readHigh / 2

	// Read Box
	drawBoxUI(readBoxmidx, readBoxmidy, editBoxWidth+2, readHigh)
	printReadmsg()

	editCurmidx := (w - editBoxWidth) / 2
	editCurmidy := h - 3
	editBoxmidx := w / 2
	editBoxmidy := h - 3
	editBoxHigh := 2

	// Edit Box
	drawBoxUI(editBoxmidx, editBoxmidy, editBoxWidth+2, editBoxHigh)
	editBox.Draw(editCurmidx, editCurmidy, editBoxWidth, 1)
	termbox.SetCursor(editCurmidx+editBox.CursorX(), editCurmidy)

	termbox.Flush()
}

// RunGui is
func RunGui(GuiInitchan chan int) {
	err := termbox.Init()
	readStrings = make([]readString, 0, 5)
	go runReadLogic()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	redrawAll()

	GuiInitchan <- 1
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				editBox.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				editBox.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				editBox.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				editBox.DeleteRuneForward()
			case termbox.KeyTab:
				editBox.InsertRune('\t')
			case termbox.KeySpace:
				editBox.InsertRune(' ')
			case termbox.KeyCtrlK:
				editBox.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				editBox.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				editBox.MoveCursorToEndOfTheLine()
			case termbox.KeyEnter:
				// Scene Handler로 메시지 전송합니다.
				Instance().ObjMgr().ChanManager().SendchanRequestFromGui(editBox.getText())
				editBox.AllClearRune()
			default:
				if ev.Ch != 0 {
					editBox.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func runeAdvanceLen(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos%tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffsetCoffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset++
		voffset += runeAdvanceLen(r, voffset)
	}
	return
}

func byteSliceGrow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 8

type EditBox struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
}

// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	const colred = termbox.ColorRed

	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, arrowRight,
				colred, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.line_voffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.line_voffset != 0 {
		termbox.SetCell(x, y, arrowLeft, colred, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.line_voffset != 0 {
		threshold = width - ht
	}
	if eb.cursor_voffset-eb.line_voffset >= threshold {
		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
	}

	if eb.line_voffset != 0 && eb.cursor_voffset-eb.line_voffset < ht {
		eb.line_voffset = eb.cursor_voffset - ht
		if eb.line_voffset < 0 {
			eb.line_voffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursor_boffset = boffset
	eb.cursor_voffset, eb.cursor_coffset = voffsetCoffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursor_boffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursor_boffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

// AllClearRune is
func (eb *EditBox) AllClearRune() {
	if eb.cursor_boffset == 0 {
		return
	}
	eb.MoveCursorToBeginningOfTheLine()
	eb.text = eb.text[:0]
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursor_boffset]
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.cursor_boffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

func (eb *EditBox) getText() string {
	return string(eb.text)
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

var arrowLeft = '←'
var arrowRight = '→'

func init() {
	if runewidth.EastAsianWidth {
		arrowLeft = '<'
		arrowRight = '>'
	}
}
