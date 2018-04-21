package tex

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/fractalqb/goxic"
)

func ExampleInlinePlaceholder() {
	tmpl := "\n\\section{`secname`}\n"
	p := NewParser()
	pts := make(map[string]*goxic.Template)
	p.Parse(strings.NewReader(tmpl), "tst", pts)
	t := pts[""]
	bt := t.NewBounT()
	bt.BindPName("secname", 4711)
	bt.Emit(os.Stdout)
	// Output:
	// \section{4711}
}

func TestBlockPlaceholder(t *testing.T) {
	txt := `\section{bla}
% >>> nested <<<
end
`
	p := NewParser()
	pts := make(map[string]*goxic.Template)
	err := p.Parse(strings.NewReader(txt), "tst", pts)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := pts[""]
	if tmpl.FixCount() != 2 {
		t.Fatal("only one fixed part")
	}
	if tmpl.PlaceholderAt(1) != "nested" {
		t.Fatalf("wrong placeholder: %s", tmpl.PlaceholderAt(0))
	}
}

func ExampleNetedTemplate() {
	tmpl := `\section{bla}
% >>> nested >>>
% <<< nested <<<
end
`
	p := NewParser()
	pts := make(map[string]*goxic.Template)
	err := p.Parse(strings.NewReader(tmpl), "tst", pts)
	if err != nil {
		fmt.Println(err)
	}
	t := pts[""]
	bt := t.NewBounT()
	bt.Emit(os.Stdout)
	// Output:
	// \section{bla}
	//
	// end
}
