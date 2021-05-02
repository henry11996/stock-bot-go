package main

import (
	"bytes"
	"html/template"

	"github.com/RainrainWu/fugle-realtime-go/client"
)

func GetInfo(data client.FugleAPIData) string {
	t := template.New("info.html")

	var err error
	t, err = t.ParseFiles("assets/html/info.html")
	if err != nil {
		return ""
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return ""
	}

	return tpl.String()
}
