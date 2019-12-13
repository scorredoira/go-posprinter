package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

func listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printHandler(w, r)
	})

	if err := http.ListenAndServe(":9823", nil); err != nil {
		log.Fatal(err)
	}
}

func printHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "auth,clientguid,pragma,Content-Type,User-Agent,Referer,Origin")
	w.Header().Set("Access-Control-Expose-Headers", "Cache-Breaker")

	data := r.FormValue("data")

	if data == "" {
		data = r.PostFormValue("data")
	}

	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Decode error", err)
		w.Write([]byte(err.Error()))
	}

	s, err := url.QueryUnescape(string(b))
	if err != nil {
		fmt.Println("Encode error", err)
		w.Write([]byte(err.Error()))
	}

	//b = []byte(s)

	b, err = charmap.Windows1252.NewEncoder().Bytes([]byte(s))
	if err != nil {
		// TODO validate chars and replace invalid with ?
		// For now just print the original text.
		fmt.Println("ENCODING", err)
	}

	if err := printer.Write(b); err != nil {
		fmt.Println("Write error", err)
		w.Write([]byte(err.Error()))
	}
}

type Printer struct {
	CUPS bool
	Name string
	file io.Writer
}

func New(name string, cups bool) (*Printer, error) {
	if cups {
		return &Printer{Name: name, CUPS: true}, nil
	}

	if strings.HasPrefix(name, "\\") {
		return &Printer{Name: name}, nil
	}

	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return &Printer{Name: name, file: f}, nil
}

func (p *Printer) sendToCUPS(b []byte) error {
	cmd := exec.Command("lpr", "-o", "raw", "-H", "localhost", "-P", p.Name)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		io.Copy(pipe, bytes.NewReader(b))
		pipe.Close()
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

func decode(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	for i, step, l := 0, 0, len(b); i < l; i += step {
		r, width := utf8.DecodeRune(b[i:])

		// \b marks the start of a byte sequence
		if r == '\b' {
			// the next byte is the length of the sequence
			s := string(b[i+1 : i+3])
			step = 3
			// Use 00 for length greater than 255 in the next 4 bytes
			if s == "00" {
				s = string(b[i+3 : i+11])
				step = 11
			}

			// convert the hex string to the number of bytes
			h, err := strconv.ParseInt(s, 16, 0)
			if err != nil {
				return nil, err
			}

			// a byte is expressed as 2 hex digits
			numBytes := int(h) * 2

			// decode hex digits to bytes
			dst := make([]byte, int(h))
			_, err = hex.Decode(dst, b[i+step:i+step+numBytes])
			if err != nil {
				return nil, fmt.Errorf("Invalid hex value (%s): %v", b[i+step:i+step+numBytes], err)
			}
			buf.Write(dst)

			step += numBytes
			continue
		}

		buf.Write(b[i : i+width])
		step = width
	}

	return buf.Bytes(), nil
}
