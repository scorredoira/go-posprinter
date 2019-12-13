package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"scorredoira/ini"
	"strconv"
)

var printer *Printer

func main() {
	fmt.Println("Point Of Sale Ticket Printer")
	fmt.Println("Version 1.03 - Amura Systems")
	fmt.Println()

	var path string
	var cups bool
	flag.StringVar(&path, "p", "", "The printer name or path (/dev/usb/lp0)")
	flag.BoolVar(&cups, "cups", false, "Use CUPS")
	flag.Parse()

	c, err := ini.New("posprinter.conf")
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
