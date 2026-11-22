package gui

import (
	"strings"
	"unicode"
)

// Answer is one button of a window in the middle: its letter, its label and what it does.
type Answer struct {
	Key   rune
	Label string
	Do    func() // nil = only close the window
}

// Dialog is a window in the middle with a question and its answers, answered only with the keyboard:
// a letter, or ← → to choose and Enter. Esc gives the answer at index Esc. Other keys are ignored.
type Dialog struct {
	Question string
	Body     []Line // optional lines under the question (e.g. where files go and their names)
	Answers  []Answer
	Chosen   int // the answer Enter gives; set when it opens
	Esc      int // the answer Esc gives
}

// YesNo is the Yes / No window. yes says whether Yes is chosen when it opens. Esc is No.
func YesNo(question string, onYes func(), yes bool) *Dialog {
	c := &Dialog{Question: question, Answers: []Answer{{Key: 'y', Label: "Yes", Do: onYes}, {Key: 'n', Label: "No"}}, Esc: 1, Chosen: 1}
	if yes {
		c.Chosen = 0
	}
	return c
}

// ChosenKey is the letter of the answer chosen now.
func (c *Dialog) ChosenKey() rune { return c.Answers[c.Chosen].Key }

func (c *Dialog) give(i int) {
	if do := c.Answers[i].Do; do != nil {
		do()
	}
}

// Key handles a key while the window is open; it returns true when the window closes.
func (c *Dialog) Key(k string) bool {
	switch k {
	case KeyLeft:
		c.Chosen = max(c.Chosen-1, 0)
	case KeyRight:
		c.Chosen = min(c.Chosen+1, len(c.Answers)-1)
	case KeyEnter:
		c.give(c.Chosen)
		return true
	case KeyEsc:
		c.give(c.Esc)
		return true
	}
	return false
}

// Rune handles a letter; it returns true when the window closes.
func (c *Dialog) Rune(r rune) bool {
	for i, a := range c.Answers {
		if unicode.ToLower(r) == a.Key {
			c.give(i)
			return true
		}
	}
	return false
}

// Hint is the hint line while the window is open, e.g. "y yes · n no · ←→ choose · enter confirm · esc no".
func (c *Dialog) Hint() string {
	var parts []string
	for _, a := range c.Answers {
		parts = append(parts, string(a.Key)+" "+strings.ToLower(a.Label))
	}
	parts = append(parts, "←→ choose", "enter confirm", "esc "+strings.ToLower(c.Answers[c.Esc].Label))
	return strings.Join(parts, " · ")
}

// Box draws the window: a frame around the question and the answers.
func (c *Dialog) Box() []Line {
	var answers Line
	for i, a := range c.Answers {
		if i > 0 {
			answers = append(answers, Plain("   "))
		}
		s := Strong(" ["+strings.ToUpper(string(a.Key))+"] "+a.Label+" ", White)
		if i == c.Chosen {
			s.Bg = SelectedBg
			s.Color = Accent
		}
		answers = append(answers, s)
	}
	inner := max(len([]rune(c.Question)), LineLen(answers))
	for _, l := range c.Body {
		inner = max(inner, LineLen(l))
	}
	inner += 6

	frame := func(l Line) Line {
		out := Line{Colored("│", Muted), Plain("   ")}
		out = append(out, l...)
		out = append(out, Plain(strings.Repeat(" ", inner-3-LineLen(l))), Colored("│", Muted))
		return out
	}
	box := []Line{
		{Colored("╭"+strings.Repeat("─", inner)+"╮", Muted)},
		frame(nil),
		frame(Line{Strong(c.Question, White)}),
		frame(nil),
	}
	if len(c.Body) > 0 {
		for _, l := range c.Body {
			box = append(box, frame(l))
		}
		box = append(box, frame(nil))
	}
	return append(box,
		frame(answers),
		frame(nil),
		Line{Colored("╰"+strings.Repeat("─", inner)+"╯", Muted)},
	)
}

// LineLen is how many characters the line holds.
func LineLen(l Line) int {
	n := 0
	for _, s := range l {
		n += len([]rune(s.Text))
	}
	return n
}

// Overlay puts box over lines, centered in a console cols wide, between the rows first and last.
func Overlay(lines []Line, box []Line, cols, first, last int) {
	if len(box) == 0 {
		return
	}
	row := first + max((last-first-len(box))/2, 0)
	col := max((cols-LineLen(box[0]))/2, 0)
	for i, b := range box {
		if r := row + i; r < len(lines) {
			lines[r] = place(lines[r], b, col)
		}
	}
}

// place writes over the characters of l from column col on, keeping the colors of what is left.
func place(l, over Line, col int) Line {
	type cell struct {
		r rune
		s Span
	}
	var cells []cell
	for _, s := range l {
		for _, r := range s.Text {
			cells = append(cells, cell{r, s})
		}
	}
	for len(cells) < col {
		cells = append(cells, cell{' ', Span{}})
	}
	end := col + LineLen(over)
	var out Line
	add := func(r rune, s Span) {
		if n := len(out); n > 0 && out[n-1].Color == s.Color && out[n-1].Bg == s.Bg && out[n-1].Bold == s.Bold {
			out[n-1].Text += string(r)
			return
		}
		s.Text = string(r)
		out = append(out, s)
	}
	for _, c := range cells[:col] {
		add(c.r, c.s)
	}
	out = append(out, over...)
	if end < len(cells) {
		for _, c := range cells[end:] {
			add(c.r, c.s)
		}
	}
	return out
}

// WithDialog draws d over the lines of a screen that ends with its hint line (cols x rows characters):
// the hint becomes the dialog's own, at the bottom of the window, and the box sits in the middle.
func WithDialog(lines []Line, d *Dialog, cols, rows int) []Line {
	out := append([]Line(nil), lines...)
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	for len(out) < rows-2 { // Draw leaves the first row empty
		out = append(out, nil)
	}
	Overlay(out, d.Box(), cols, 0, len(out))
	return append(out, Line{Colored(d.Hint(), Muted)})
}

// QuitDialog asks before closing an app, with Yes chosen (author, 2026-11-19: the same in every app).
func QuitDialog(name string, quit func()) *Dialog {
	return YesNo("Quit "+name+"?", quit, true)
}
