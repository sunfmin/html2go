# Generate html to htmlgo code

Install

```bash
$ go get github.com/sunfmin/html2go
```

Run 

```bash
$ html2go # enter and then copy paste html code
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <input readonly required disabled checked tabindex="-1">
  <input readonly="false">
</nav>
# enter a new line here
```

Then ctrl+d, the terminal will show:

```bash
package hello

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
```

Use `-pkg` if you need to generate a different package prefix

```bash
$ html2go -pkg=h
```
