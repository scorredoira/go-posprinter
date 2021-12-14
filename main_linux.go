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
	fmt.Println("Version 1.03 - Amura Systems")
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

	// if err := printer.Write([]byte("0x210x40")); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := printer.Write([]byte("hola hola\n\n")); err != nil {
	// 	log.Fatal(err)
	// }

	content_bytes := []byte("A:123456789*B:999999990*C:PT*D:FR*E:N*F:20211207*G:FTR 121/81*H:0*I1:PT*I7:5.69*I8:1.31*N:1.31*O:7.00*Q:Nhz9*R:0181")
	store_len := len(content_bytes) + 3
	store_pL := (byte)(store_len % 256)
	store_pH := (byte)(store_len / 256)

	INIT := []byte{0x1B, 0x40}
	FUNC_165 := []byte{0x1D, 0x28, 0x6b, 0x04, 0x00, 0x31, 0x41, 0x33, 0x00}
	FUNC_167 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x43, 0x05} // el último es el tamaño 3-10
	FUNC_169 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x45, 0x30}
	FUNC_180 := []byte{0x1D, 0x28, 0x6b, store_pL, store_pH, 0x31, 0x50, 0x30}
	FUNC_181 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x51, 0x30}
	FUNC_182 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x52, 0x30}

	printer.Write(INIT)
	printer.Write(FUNC_165)
	printer.Write(FUNC_167)
	printer.Write(FUNC_169)
	printer.Write(FUNC_180)
	printer.Write(content_bytes)
	printer.Write(FUNC_181)
	printer.Write(FUNC_182)

	// if err := printer.Write([]byte{})); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := printer.Write([]byte("\x1d\x6b\x04")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("HELLO")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x00")); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := printer.Write([]byte("\x1D\x28\x6B\x04\x00\x31\x41\x33\x00")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x04\x00\x31\x41\x33\x00")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x04\x00\x31\x41\x33\x00")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x03\x00\x31\x43\x03")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x03\x00\x31\x45\x30")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x34\x00\x31\x50\x30")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x74\x65\x73\x74")); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printer.Write([]byte("\x1D\x28\x6B\x03\x00\x31\x51\x48")); err != nil {
	// 	log.Fatal(err)
	// }

	return

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
