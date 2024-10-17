package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/zepyrshut/gorender"
)

func dummyFunc() string {
	return "dummy function"
}

func main() {
	newFuncs := template.FuncMap{
		"dummyFunc": dummyFunc,
	}

	renderOpts := &gorender.Render{
		EnableCache:       true,
		TemplatesPath:     "template",
		PageTemplatesPath: "template/pages",
		Functions:         newFuncs,
	}

	ren := gorender.New(gorender.WithRenderOptions(renderOpts))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		td := &gorender.TemplateData{}

		ren.Template(w, r, "page.html", td)
	})

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
