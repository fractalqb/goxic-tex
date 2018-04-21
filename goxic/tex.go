package tex

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"unicode/utf8"

	"github.com/fractalqb/goxic"
)

func NewParser() *goxic.Parser {
	res := goxic.Parser{
		StartInlinePh: "`",
		EndInlinePh:   "`",
		BlockPh: regexp.MustCompile(
			`^[ \t]*%(\\?) >>> ([a-zA-Z0-9_-]+) <<<[ \t]*(\\?)[ \t]*$`),
		PhNameRgxGrp: 2,
		PhLBrkRgxGrp: 1,
		PhTBrkRgxGrp: 3,
		StartSubTemplate: regexp.MustCompile(
			`^[ \t]*%(\\?) >>> ([a-zA-Z0-9_-]+) >>>[ \t]*$`),
		StartNameRgxGrp: 2,
		StartLBrkRgxGrp: 1,
		EndSubTemplate: regexp.MustCompile(
			`^[ \t]*% <<< ([a-zA-Z0-9_-]+) <<<[ \t]*(\\?)[ \t]*$`),
		EndNameRgxGrp: 1,
		EndTBrkRgxGrp: 2,
		Endl:          "\n"}
	return &res
}

type EscWriter struct {
	Escape io.Writer
	buf    [utf8.UTFMax]byte
	wp     int
}

func (hew *EscWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		hew.buf[hew.wp] = b
		hew.wp++
		if buf := hew.buf[:hew.wp]; utf8.FullRune(buf) {
			hew.wp = 0
			if r, _ := utf8.DecodeRune(buf); r == utf8.RuneError {
				return n, errors.New("utf8 rune decoding error")
			} else {
				switch r {
				case '\\':
					if i, err := hew.Escape.Write([]byte(`\\`)); err != nil {
						return n + i, err
					} else {
						n += i
					}
				case '%':
					if i, err := hew.Escape.Write([]byte(`\%`)); err != nil {
						return n + i, err
					} else {
						n += i
					}
				case '&':
					if i, err := hew.Escape.Write([]byte(`\&`)); err != nil {
						return n + i, err
					} else {
						n += i
					}
				case '$':
					if i, err := hew.Escape.Write([]byte(`\$`)); err != nil {
						return n + i, err
					} else {
						n += i
					}
				default:
					if i, err := hew.Escape.Write([]byte(string(r))); err != nil {
						return n + i, err
					} else {
						n += i
					}
				}
			}
		}
	}
	return n, nil
}

func Escape(str string) string {
	buf := bytes.NewBuffer(nil)
	ewr := EscWriter{Escape: buf}
	if _, err := ewr.Write([]byte(str)); err != nil {
		panic(err)
	}
	return buf.String()
}

type Escaper struct {
	Cnt goxic.Content
}

func (hc Escaper) Emit(wr io.Writer) int {
	esc := EscWriter{Escape: wr}
	return hc.Cnt.Emit(&esc)
}
