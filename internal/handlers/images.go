package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
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

var seq = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// imageUploadResult hold image upload data
type imageUploadResult struct {
	Url   string
	Error error
}

// parseImageTags parses image tags and returns correct comment content
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
		return "", err
	}
	var buf bytes.Buffer
	ch := make(chan imageUploadResult, 10)
	for _, node := range content {
		if node.Type == html.ElementNode && node.Data == "img" {
			for i, img := range node.Attr {
				if img.Key == "src" {
					go repo.saveImages(img.Val, ch)
					url := <-ch
					if url.Error != nil {
						return "", url.Error
					}
					node.Attr[i].Val = url.Url
					break
				}
			}
		}
		if err := html.Render(&buf, node); err != nil {
			return "", err
		}
	}
	close(ch)
	return buf.String(), nil
}

// saveImages saves base64 images into the server and returns the full url to images
func (repo *Repository) saveImages(data string, ch chan imageUploadResult) {
	resp := imageUploadResult{Url: "", Error: nil}
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		resp.Error = errors.New("invalid image")
		ch <- resp
	}
	ImageType := data[11:idx]
	rand.Seed(time.Now().UnixNano())
	imgName := repo.App.StaticImages + randomSequence(15)
	var url string
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		resp.Error = errors.New("cannot decode base64 image")
		ch <- resp
	}
	r := bytes.NewReader(unbased)
	switch ImageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode png")
			ch <- resp
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.png", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			resp.Error = errors.New("cannot open file")
			ch <- resp
		}
		err = png.Encode(f, im)
		if err != nil {
			resp.Error = errors.New("cannot encode png file")
			ch <- resp
		}
		url = f.Name()

	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode jpeg")
			ch <- resp
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.jpeg", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			resp.Error = errors.New("cannot decode jpeg file")
			ch <- resp
		}
		err = jpeg.Encode(f, im, nil)
		if err != nil {
			resp.Error = errors.New("cannot encode jpeg file")
			ch <- resp
		}
		url = f.Name()

	case "gif":
		im, err := gif.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode gif")
			ch <- resp
		}
		f, err := os.OpenFile(fmt.Sprintf("%s.gif", imgName), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			resp.Error = errors.New("cannot open gif file")
			ch <- resp
		}
		err = gif.Encode(f, im, nil)
		if err != nil {
			resp.Error = errors.New("cannot encode gif")
			ch <- resp
		}
		url = f.Name()
	}
	fullUrl := fmt.Sprintf("http://%s/%s", os.Getenv("APP_BASE_URL"), url)
	resp.Url = fullUrl
	ch <- resp
}

// randomSequence returns random sequence from the provided sequence
func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = seq[rand.Intn(len(seq))]
	}
	return string(b)
}
