package service

import (
	"fmt"
	"io/ioutil"
)

func HtmlContent(htmltext string) []byte {
	htmlContent, err := ioutil.ReadFile(htmltext)
	if err != nil {
		fmt.Println("Ошибка чтения HTML-файла", err)
		return nil
	}
	return htmlContent
}
