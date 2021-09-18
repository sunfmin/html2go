package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/theplant/htmlgo"
	"golang.org/x/net/html"
)

func GenerateHTMLGo(pkg string, htmlCode io.Reader) string {
	n, err := html.Parse(htmlCode)
	if err != nil {
		panic(err)
	}
	methodNames := tagMethodNames()
	fc := &funcCall{}
	walk(n.FirstChild.FirstChild.NextSibling, fc, methodNames)

	code := string(fc.MarshalCode(methodNames, pkg))
	code = strings.TrimRight(code, ",\n")

	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, "", "package hello\n var n = "+code, 0)
	if err != nil {
		hl, _ := strconv.ParseInt(strings.Split(err.Error(), ":")[0], 10, 64)
		panic(fmt.Sprintf("%s\n%s", err, codeWithLineNumber(code, hl)))
	}
	buf := bytes.NewBuffer(nil)
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func codeWithLineNumber(src string, highlightLine int64) (r string) {
	lines := strings.Split(src, "\n")
	linesWithNumber := []string{}
	for i, l := range lines {
		hl := "   "
		if int64(i+1) == highlightLine {
			hl = ">> "
		}
		linesWithNumber = append(linesWithNumber, fmt.Sprintf("%s%d: %s", hl, i+1, l))
	}
	r = strings.Join(linesWithNumber, "\n")
	return
}

func pkgDot(pkg string) (r string) {
	if len(pkg) == 0 {
		return
	}
	return pkg + "."
}

type funcCall struct {
	Pkg      string
	Name     string
	Text     string
	TakeText bool
	Children []*funcCall
	Attrs    []html.Attribute
}

func (fc *funcCall) MarshalCode(methodNames []string, pkg string) (r []byte) {

	buf := bytes.NewBuffer(nil)

	if len(fc.Text) > 0 {
		buf.WriteString(fmt.Sprintf("%sText(%#+v),\n", pkgDot(pkg), fc.Text))
		return buf.Bytes()
	}

	newline := "\n"
	if fc.TakeText {
		newline = ""
	}
	_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(pkg), strcase.ToCamel(fc.Name), newline)

	needWriteChilren := false
	if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
		buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
	} else if fc.TakeText {
		buf.WriteString(`""`)
		needWriteChilren = true
	} else {
		for _, c := range fc.Children {
			buf.Write(c.MarshalCode(methodNames, pkg))
		}
	}

	buf.WriteString(")")
	for i, att := range fc.Attrs {
		attFuncName := getFuncName(att.Key, methodNames)

		buf.WriteString(".")
		if i > 0 {
			buf.WriteString("\n")
		}

		if len(attFuncName) > 0 {
			var val interface{} = att.Val
			if strings.Index(boolAttr, "|"+attFuncName+"|") >= 0 {
				val = true
			}
			if strings.Index(intAttr, "|"+attFuncName+"|") >= 0 {
				var err error
				val, err = strconv.ParseInt(att.Val, 10, 64)
				if err != nil {
					panic(err)
				}
			}
			_, _ = fmt.Fprintf(buf, "%s(%s)", attFuncName, normalizeGoString(val))
		} else {
			_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
		}
	}

	if needWriteChilren && len(fc.Children) > 0 {
		buf.WriteString(".\n")
		buf.WriteString("Children(\n")
		for _, c := range fc.Children {
			buf.Write(c.MarshalCode(methodNames, pkg))
		}
		buf.WriteString(")")
	}

	buf.WriteString(",\n")

	return buf.Bytes()
}

func expandAlpineKey(key string) (r string) {
	if strings.Index(key, ":") == 0 {
		return fmt.Sprintf("x-bind%s", key)
	}

	if strings.Index(key, "@") == 0 {
		return fmt.Sprintf("x-on%s", key)
	}
	return key
}

func normalizeGoString(val interface{}) (r interface{}) {
	strval, ok := val.(string)
	if !ok {
		return fmt.Sprintf("%#+v", val)
	}

	if strings.Contains(strval, "'") {
		strval = strings.ReplaceAll(strval, "'", "\"")
	}
	if strings.ContainsAny(strval, "\n\t\"") && !strings.Contains(strval, "`") {
		strval = fmt.Sprintf("`%s`", strval)
		return strval
	}

	return fmt.Sprintf("%#+v", strval)
}

const intAttr = "|TabIndex|"
const boolAttr = "|Required|Readonly|Disabled|Checked|"

const textTags = "|Abbr|B|Bdi|Bdo|Button|Caption|Code|Del|Dfn|Em|Figcaption|H1|H2|H3|H4|H5|" +
	"H6|I|Img|Input|Kbd|Label|Legend|Link|Mark|Object|Option|Param|Pre|Q|Rp|Rt|S|" +
	"Script|Small|Source|Span|Strong|Style|Sub|Sup|Textarea|Th|Time|Title|Track|U|Var|Wbr|"

func walk(n *html.Node, fc *funcCall, methodNames []string) {
	switch n.Type {
	case html.ElementNode:
		if len(strings.TrimSpace(n.Data)) > 0 {
			fc.Name = strcase.ToCamel(strings.TrimSpace(n.Data))
		}
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) > 0 {
			fc.Text = strings.TrimSpace(n.Data)
		}
	}

	if strings.Index(textTags, "|"+fc.Name+"|") >= 0 {
		fc.TakeText = true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && len(strings.TrimSpace(c.Data)) == 0 {
			continue
		}
		if c.Type == html.CommentNode {
			continue
		}

		ch := &funcCall{Attrs: c.Attr}
		fc.Children = append(fc.Children, ch)
		walk(c, ch, methodNames)
	}
}

func getFuncName(name string, methodNames []string) (r string) {
	for _, m := range methodNames {
		if strings.ToLower(name) == strings.ToLower(m) {
			return m
		}
	}
	return ""
}

func tagMethodNames() (r []string) {
	tag := htmlgo.Tag("")
	tagType := reflect.TypeOf(tag)
	for i := 0; i < tagType.NumMethod(); i++ {
		r = append(r, tagType.Method(i).Name)
	}
	return
}
