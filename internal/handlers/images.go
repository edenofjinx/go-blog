package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (repo *Repository) parseImageTags(s string) (string, error) {
	content, err := html.ParseFragment(
		strings.NewReader(s),
		&html.Node{
			Type:     html.ElementNode,
			Data:     "body",
			DataAtom: atom.Body,
		},
	)
	if err != nil {
		log.Println(err)
	}
	var buf bytes.Buffer
	for _, node := range content {
		if node.Type == html.ElementNode && node.Data == "img" {
			for i, img := range node.Attr {
				if img.Key == "src" {
					url, err := repo.saveImages(img.Val)
					if err != nil {
						return "", err
					}
					node.Attr[i].Val = url
					break
				}
			}
		}
		if err := html.Render(&buf, node); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func (repo *Repository) saveImages(data string) (string, error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		panic("InvalidImage")
	}
	ImageType := data[11:idx]
	rand.Seed(time.Now().UnixNano())
	imgName := repo.App.StaticImages + randSeq(15)
	var url string
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		panic("Cannot decode b64")
	}
	r := bytes.NewReader(unbased)
	switch ImageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			panic("Bad png")
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.png", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}
		err = png.Encode(f, im)
		if err != nil {
			return "", err
		}
		url = f.Name()

	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			panic("Bad jpeg")
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.jpeg", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}
		err = jpeg.Encode(f, im, nil)
		if err != nil {
			return "", err
		}
		url = f.Name()

	case "gif":
		im, err := gif.Decode(r)
		if err != nil {
			panic("Bad gif")
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.gif", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}
		err = gif.Encode(f, im, nil)
		if err != nil {
			return "", err
		}
		url = f.Name()
	}
	return url, nil
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
