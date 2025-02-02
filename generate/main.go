package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

var nsToURL = map[string]string{
	"atom":    "http://www.w3.org/2005/Atom",
	"itunes":  "http://www.itunes.com/dtds/podcast-1.0.dtd",
	"podcast": "https://podcastindex.org/namespace/1.0",
	"content": "http://purl.org/rss/1.0/modules/content/",
}

func main() {
	output := &bytes.Buffer{}

	fmt.Fprint(output, "// Code generated by gopodcast generator. DO NOT EDIT.\n\n")
	fmt.Fprint(output, "// This file contains copies of structs to fix an issue with namespace prefix support in go `encoding/xml` package\n\n")
	fmt.Fprint(output, "package gopodcast\n\n")
	fmt.Fprint(output, "import \"encoding/xml\"\n\n")

	fsrc, err := os.ReadFile("gopodcast.go")
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fsrc, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	type astField struct {
		name  string
		fType string
		tag   string
	}

	type astStruct struct {
		name   string
		fields []astField
	}

	structs := make([]astStruct, 0)

	lastIdent := ""
	ast.Inspect(f, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			lastIdent = id.Name
		}

		if n, ok := n.(*ast.StructType); ok {
			strct := astStruct{
				name:   lastIdent,
				fields: make([]astField, 0),
			}
			for _, f := range n.Fields.List {
				tag := ""
				if f.Tag != nil {
					tag = transformXMLTag(f.Tag.Value)
				}
				typeStr := fsrc[f.Type.Pos()-1 : f.Type.End()-1]
				strct.fields = append(strct.fields, astField{
					name:  f.Names[0].Name,
					fType: string(typeStr),
					tag:   tag,
				})
			}
			structs = append(structs, strct)
		}
		return true
	})

	for _, strct := range structs {
		fmt.Fprintf(output, "type xmlFix%s struct {\n", strct.name)
		for _, field := range strct.fields {
			fmt.Fprintf(output, "\t%s %s %s\n", field.name, newFType(field.fType), field.tag)
		}
		fmt.Fprint(output, "}\n\n")

		fmt.Fprintf(output, "func (s *xmlFix%s) Translate() *%s {\n", strct.name, strct.name)
		fmt.Fprint(output, "if s == nil {\n\treturn nil\n}\n")
		fmt.Fprintf(output, "\tvar r %s\n", strct.name)
		for _, field := range strct.fields {
			if isIgnoreType(field.fType) {
				fmt.Fprintf(output, "\tr.%s = s.%s\n", field.name, field.name)
				continue
			}

			if isPtr(field.fType) {
				fmt.Fprintf(output, "\tr.%s = s.%s.Translate()\n", field.name, field.name)
			} else if isSlicePtr(field.fType) {
				fmt.Fprintf(output, "\tv%s := make(%s, 0)\n", field.name, field.fType)
				fmt.Fprintf(output, "\tfor _, v := range s.%s {\n", field.name)
				fmt.Fprintf(output, "\t\tv%s = append(v%s, v.Translate())\n", field.name, field.name)
				fmt.Fprint(output, "\t}\n")
				fmt.Fprintf(output, "\tr.%s = v%s\n", field.name, field.name)
			} else if isSlice(field.fType) {
				fmt.Fprintf(output, "\tv%s := make(%s, 0)\n", field.name, field.fType)
				fmt.Fprintf(output, "\tfor _, v := range s.%s {\n", field.name)
				fmt.Fprint(output, "\t\tx := v.Translate()\n")
				fmt.Fprintf(output, "\t\tv%s = append(v%s, *x)\n", field.name, field.name)
				fmt.Fprintf(output, "\t}\n")
				fmt.Fprintf(output, "\tr.%s = v%s\n", field.name, field.name)
			} else {
				fmt.Fprintf(output, "\tv%s := s.%s.Translate()\n", field.name, field.name)
				fmt.Fprintf(output, "\tr.%s = *v%s\n", field.name, field.name)
			}
		}
		fmt.Fprint(output, "\treturn &r\n")
		fmt.Fprint(output, "}\n\n")
	}

	src, err := format.Source(output.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("gopodcast_xml_fix.go", src, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func transformXMLTag(s string) string {
	r := regexp.MustCompile(`xml:"(([a-zA-Z]+):([a-zA-Z]+))[,"].*`)
	matches := r.FindStringSubmatch(s)
	if len(matches) != 4 {
		return s
	}

	// use NS URL instead of NS when parsing
	ns := matches[2]
	if nsURL, ok := nsToURL[ns]; ok {
		ns = nsURL
	}

	return strings.ReplaceAll(s, matches[1], fmt.Sprintf("%s %s", ns, matches[3]))
}

var ignoreTypes = []string{"string", "bool", "int", "int64", "byte", "xml.Name", "FlexBool"}

func newFType(inType string) string {
	prefix := ""
	fType := inType
	n := strings.LastIndexAny(inType, "*[]")
	if n >= 0 {
		prefix = inType[:n+1]
		fType = inType[n+1:]

	}
	if slices.Contains(ignoreTypes, fType) {
		fType = prefix + fType
	} else {
		fType = prefix + "xmlFix" + fType
	}
	return fType
}

func isIgnoreType(t string) bool {
	if slices.Contains(ignoreTypes, strings.TrimLeft(t, "*[]")) {
		return true
	}
	return false
}

func isPtr(t string) bool {
	return strings.HasPrefix(t, "*")
}

func isSlice(t string) bool {
	return strings.HasPrefix(t, "[")
}

func isSlicePtr(t string) bool {
	return strings.HasPrefix(t, "[]*")
}
