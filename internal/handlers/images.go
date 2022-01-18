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
	"path/filepath"
	"runtime"
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
		return
	}
	imageType := data[11:idx]
	rand.Seed(time.Now().UnixNano())
	imgName := randomSequence(15)
	var url string
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		resp.Error = errors.New("cannot decode base64 image")
		ch <- resp
		return
	}
	r := bytes.NewReader(unbased)
	switch imageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode png")
			ch <- resp
			return
		}
		f, err := repo.openImageFile(imageType, imgName)
		if err != nil {
			resp.Error = err
			ch <- resp
			return
		}
		err = png.Encode(f, im)
		if err != nil {
			resp.Error = errors.New("cannot encode png file")
			ch <- resp
			return
		}
		url = f.Name()

	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode jpeg")
			ch <- resp
			return
		}
		f, err := repo.openImageFile(imageType, imgName)
		if err != nil {
			resp.Error = err
			ch <- resp
			return
		}
		err = jpeg.Encode(f, im, nil)
		if err != nil {
			resp.Error = errors.New("cannot encode jpeg file")
			ch <- resp
			return
		}
		url = f.Name()

	case "gif":
		im, err := gif.Decode(r)
		if err != nil {
			resp.Error = errors.New("cannot decode gif")
			ch <- resp
			return
		}
		f, err := repo.openImageFile(imageType, imgName)
		if err != nil {
			resp.Error = err
			ch <- resp
			return
		}
		err = gif.Encode(f, im, nil)
		if err != nil {
			resp.Error = errors.New("cannot encode gif")
			ch <- resp
			return
		}
		url = f.Name()
	default:
		resp.Error = errors.New(fmt.Sprintf("incorrect image type provided: %s", imageType))
		ch <- resp
		return
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

func (repo *Repository) openImageFile(t string, n string) (*os.File, error) {
	path, err := repo.getStaticImageDir()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s%s.%s", path, n, t), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot open file type %s", t))
	}
	return f, nil
}

func (repo *Repository) getStaticImageDir() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	imgPath := root + "/" + repo.App.StaticImages
	errDir := os.MkdirAll(imgPath, 0755)
	if errDir != nil {
		return "", errors.New("could not create directory")
	}
	return imgPath, nil
}
