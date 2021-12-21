package main

import (
	"encoding/base64"
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
	// FUNC_167 := []byte{0x1D, 0x28, 0x6b, 0x03, 0x00, 0x31, 0x43, 0x09} // el último es el tamaño 3-9
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

	b, err := base64.RawStdEncoding.DecodeString("G3QDR29sZiBEZW1vICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIA0KUE9SVFVHQUwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIA0KMDAwMC0wMDAgUE9SVFVHQUwgICAgICAgICAgICAgICAgICAgICAgDQpDb250cmliLjogMTIzNDU2Nzg5ICAgICAgICAgICANClRlbDogICAgICAgICAgICAgICAgIEZheDogICAgICAgICAgICAgICANCg0KDQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tDQpPcGVyYWRvcjogZ29sZiAgICAgICAgICAgICAgICANCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0NCg0KQ29kaWdvOiAgICAgMQ0KTm9tZS4uOiBJbmVzIFBhcmVudGUgICAgICAgICAgICAgICAgICAgIA0KTW9yYWRhOiBibGFibGEgICAgICAgICAgICAgICAgICAgICAgICAgIA0KQy5Qb3N0YWw6bnVsbCAgICAgUFQgICAgICAgICAgICAgICAgICANCkNvbnRyaWIuOiBFUzU1ODk2NjYgICAgICAgICAgIA0KDQo9PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09DQpGYXR1cmEgLSBSZWNpYm8gICAgICBOLkZUUiAxMjEvOTAgICAgDQoyMDIxLTEyLTE2ICAgICAgICAgICAgICAgICAgICAgICAgIDIwOjE0DQo9PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09DQogIFF0ZCBBcnRpZ28gICAgICAgICAgICAgICAgJUlWQSAgIFRvdGFsDQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tDQoqKioqKiBQcm9kdXRvIEVtcHJlc2EgOTk4ICAgIDIzICAgIDIwLDAwDQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tDQogICAgVG90YWwuLi4uLi4uLi4uLi4uLi4uLi46ICAgICAgIDIwLDAwDQogICAgLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tDQogICAgQ2hlcXVlICAgICAgICAgICAgICAgICAgICAgICAgIDIwLDAwDQoNCiAgICBUcm9jby4uLi4uLi4uLi4uLi4uLi4uLjogICAgICAgIDAsMDANCiAgICAtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0NCiAgICBUb3RhbCBkZSBEZXNjb250b3MuLi4uLjogICAgICAgIDAsMDANCiAgICAtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0NCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICANCj09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT0NCkRlc2NyaWNhbyAgICAgVGF4YSAgIEluY2lkZW5jaWEgICAgVmFsb3INCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0NCklWQSAgICAgICAgICAgIDIzJSAgICAgIDE2LDI2ICAgICAgIDMsNzQNCj09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT0NCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0NCg0KG2ExHShrBAAxQTIAHShrAwAxQwgdKGsDADFFMR0oa3YAMVAwQToxMjM0NTY3ODkqQjo1NTg5NjY2KkM6UFQqRDpGUipFOk4qRjoyMDIxMTIxNipHOkZUUiAxMjEvOTAqSDowKkkxOlBUKkk3OjE2LjI2Kkk4OjMuNzQqTjozLjc0Kk86MjAuMDAqUTpBWnQxKlI6MDE4MR0oawMAMVEwG2EwT0JSSUdBRE8sIFZPTFRFIFNFTVBSRSENCmV0aWNhZGF0YSAtIEFadDEtUFJPQ0VTU0FETyBQT1IgUFJPR1JBTUENCiAgICAgICAgIENFUlRJRklDQURPIE4uMDE4MS9BVCAgICAgICAgICANCk9yaWdpbmFsICAgICAgICAgICAgDQoNCg0KDQoNCg0KDQoNCg0KDQoNCg==")

	// const N = 20
	// for i, len := 0, len(b); i < len; i += N {
	// 	if i+N >= len {
	// 		if err := printer.Write(b[i:]); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	} else {
	// 		if err := printer.Write(b[i : i+N]); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

	for _, c := range b {
		if err := printer.Write([]byte{c}); err != nil {
			log.Fatal(err)
		}
	}

	//	fmt.Println(string(b))

	//printer.Write(b)
	return

	listen()
}

func usage() {
	fmt.Println(`Usage: posprinter -p /dev/usb/lpt1`)
	os.Exit(1)
}

func (p *Printer) Write(b []byte) error {
	// d, err := decode(b)
	// if err != nil {
	// 	return err
	// }

	d := b

	if p.CUPS {
		return p.sendToCUPS(d)
	}

	if _, err := p.file.Write(d); err != nil {
		return err
	}

	return nil
}
