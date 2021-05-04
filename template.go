package main

import (
	"bytes"
	"html/template"

	"github.com/RainrainWu/fugle-realtime-go/client"
)

func convertByTemplate(templa string, data client.FugleAPIData) (string, error) {
	t := template.New(templa + ".html")

	var err error
	t, err = t.ParseFiles("assets/html/" + templa + ".html")
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
