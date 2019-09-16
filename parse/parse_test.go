package parse_test

import (
	"strings"
	"testing"

	"github.com/sunfmin/html2go/parse"
	"github.com/theplant/testingutils"
)

func TestAll(t *testing.T) {
	var cases = []struct {
		name   string
		pkg    string
		html   string
		gocode string
	}{
		{
			name: "normal",
			html: `
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <a class="navbar-brand" href="#">Navbar</a>
  <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbarNav">
    <ul class="navbar-nav">
      <li class="nav-item active">
        <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="#">Features</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="#">Pricing</a>
      </li>
      <li class="nav-item">
        <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>
      </li>
    </ul>
  </div>
</nav>
`,
			gocode: `package hello

var n = Body(
	Nav(
		A(
			Text("Navbar"),
		).Class("navbar-brand").
			Href("#"),
		Button("").Class("navbar-toggler").
			Type("button").
			Attr("data-toggle", "collapse").
			Attr("data-target", "#navbarNav").
			Attr("aria-controls", "navbarNav").
			Attr("aria-expanded", "false").
			Attr("aria-label", "Toggle navigation").
			Children(
				Span("").Class("navbar-toggler-icon"),
			),
		Div(
			Ul(
				Li(
					A(
						Text("Home"),
						Span("(current)").Class("sr-only"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item active"),
				Li(
					A(
						Text("Features"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item"),
				Li(
					A(
						Text("Pricing"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item"),
				Li(
					A(
						Text("Disabled"),
					).Class("nav-link disabled").
						Href("#").
						TabIndex(-1).
						Attr("aria-disabled", "true"),
				).Class("nav-item"),
			).Class("navbar-nav"),
		).Class("collapse navbar-collapse").
			Id("navbarNav"),
	).Class("navbar navbar-expand-lg navbar-light bg-light"),
)
`,
		},
		{
			name: "bool attribute always be true even if value is false",
			html: `
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <input readonly required disabled checked tabindex="-1">
  <input readonly="false">
</nav>
`,
			gocode: `package hello

var n = Body(
	Nav(
		Input("").Readonly(true).
			Required(true).
			Disabled(true).
			Checked(true).
			TabIndex(-1),
		Input("").Readonly(true),
	).Class("navbar navbar-expand-lg navbar-light bg-light"),
)
`,
		},
		{
			name: "text attr with text",
			html: `
<div>
  <span>Hello</span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("Hello"),
	),
)
`,
		},

		{
			name: "text attr with more children",
			html: `
<div>
  <span>Hello<b>world</b></span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("").
			Children(
				Text("Hello"),
				B("world"),
			),
	),
)
`,
		},
		{
			name: "text attr on tag children",
			html: `
<div>
  <span><b>world</b></span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("").
			Children(
				B("world"),
			),
	),
)
`,
		},
		{
			name: "text attr with more children with pkg",
			pkg:  "h",
			html: `
<div>
  <span>Hello<b>world</b></span>
</div>
`,
			gocode: `package hello

var n = h.Body(
	h.Div(
		h.Span("").
			Children(
				h.Text("Hello"),
				h.B("world"),
			),
	),
)
`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gocode := parse.GenerateHTMLGo(c.pkg, strings.NewReader(c.html))
			diff := testingutils.PrettyJsonDiff(c.gocode, gocode)

			if len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}
