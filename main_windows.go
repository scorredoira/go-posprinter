package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	winPrint "github.com/alexbrainman/printer"
)

var printer *Printer

func main() {
	fmt.Println("Point Of Sale Ticket Printer")
	fmt.Println("Version 1.03.3 for Windows - Amura Systems")
	fmt.Println()

	var path string
	var cups bool
	flag.StringVar(&path, "p", "", "The printer name")
	flag.BoolVar(&cups, "cups", false, "Use CUPS")

	list := flag.Bool("l", false, "List printers")
	flag.Parse()

	if *list {
		names, err := winPrint.ReadNames()
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Printers:")
		for _, n := range names {
			fmt.Println(" - " + n)
		}

		return
	}

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

	printer = &Printer{Name: path}

	listen()
}

func usage() {
	fmt.Println(`Usage: posprinter -p name`)
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

	pr, err := winPrint.Open(p.Name)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}
	defer pr.Close()

	err = pr.StartDocument("POS_Ticket", "RAW")
	if err != nil {
		log.Fatalf("StartDocument failed: %v", err)
	}
	defer pr.EndDocument()

	err = pr.StartPage()
	if err != nil {
		log.Fatalf("StartPage failed: %v", err)
	}

	if _, err := pr.Write(d); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	err = pr.EndPage()
	if err != nil {
		log.Fatalf("EndPage failed: %v", err)
	}

	return nil
}
