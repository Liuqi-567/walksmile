package views

import (
	"html/template"
	"path/filepath"
)

var Temp *template.Template

func InitTemplate(isDebug bool, funs template.FuncMap) {
	if !isDebug {
		names := AssetNames()
		for _, v := range names {
			name := filepath.Base(v)
			data, err := Asset(v)
			if err != nil {
				panic(err)
			}

			var tmpl *template.Template
			if Temp == nil {
				Temp = template.New(name).Funcs(funs)
			}
			if name == Temp.Name() {
				tmpl = Temp
			} else {
				tmpl = Temp.New(name).Funcs(funs)
			}
			_, err = tmpl.Parse(string(data))
			if err != nil {
				panic(err)
			}
		}
		return
	}
	Temp = template.New("")
	var err error
	Temp, err = Temp.Funcs(funs).ParseGlob("views/*")
	if err != nil {
		panic(err)
	}

}
