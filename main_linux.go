package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var printer *Printer

func main() {
	fmt.Println("Point Of Sale Ticket Printer")
	fmt.Println("Version 1.04 - SCL")
	fmt.Println()

	var path string
	var cups bool
	flag.StringVar(&path, "p", "", "The printer name or path (/dev/usb/lp0)")
	flag.BoolVar(&cups, "cups", false, "Use CUPS")
	flag.Parse()

	c, err := NewConfig("posprinter.conf")
	if err != nil {
		log.Fatal(err)
	}

	if path == "" {
		path = c.String("path", "")
		cups = c.Bool("cups", false)
	} else {
		c.Set("path", path)
		c.Set("cups", strconv.FormatBool(cups))
		c.Save()
	}

	if path == "" {
		usage()
	}

	printer, err = New(path, cups)
	if err != nil {
		log.Fatal(err)
	}

	// // Print QR
	// content_bytes := []byte("Ahlkjhlkjhlkjh9i879086709871")
	// store_len := len(content_bytes) + 3
	// store_pL := (byte)(store_len % 256)
	// store_pH := (byte)(store_len / 256)

	// INIT := []byte{0x1B, 0x40}
	// FUNC_165 := []byte{0x1D, 0x28, 0x6b, 0x04, 0x00, 0x31, 0x41, 0x33, 0x00}
	// FUNC_167 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x43, 0x05} // el último es el tamaño 3-9
	// FUNC_169 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x45, 0x30}
	// FUNC_180 := []byte{0x1D, 0x28, 0x6b, store_pL, store_pH, 0x31, 0x50, 0x30}
	// FUNC_181 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x51, 0x30}
	// FUNC_182 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x52, 0x30}

	// printer.Write(INIT)
	// printer.Write(FUNC_165)
	// printer.Write(FUNC_167)
	// printer.Write(FUNC_169)
	// printer.Write(FUNC_180)
	// printer.Write(content_bytes)
	// printer.Write(FUNC_181)
	// printer.Write(FUNC_182)

	listen()
}

func usage() {
	fmt.Println(`Usage: posprinter -p /dev/usb/lpt1`)
	os.Exit(1)
}

func (p *Printer) Write(b []byte) error {
	d, err := decode(b)
	if err != nil {
		return err
	}

	if p.CUPS {
		return p.sendToCUPS(d)
	}

	if _, err := p.file.Write(d); err != nil {
		return err
	}

	return nil
}
