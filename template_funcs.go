package main

import (
	"html/template"
)

var FuncMap = template.FuncMap{
	"inc": Inc,
}

func Inc(i int) int {
	return i + 1
}
