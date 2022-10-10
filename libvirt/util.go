package libvirt

import (
	"bufio"
	"fmt"
	"log"
	"os"

	etree "github.com/beevik/etree"
)

// TODO : get XML and save as file then read from file ?

// func Get_xml_path(xml string, path string, function func()) {
// 	var doc string
// 	var result string
// }

func Get_xpath(doc *etree.Element, path string) *etree.Element {
	var result *etree.Element
	ret := doc.FindElement(path)
	if len(ret.Child) >= 1 {
		result = ret.ChildElements()[0]
	} else {
		result = ret
	}
	return result
}

func ReadWrite(input string, output_file string) {
	f, err := os.Create(output_file)
	check_panic(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	n, err := w.WriteString(input)
	check_panic(err)
	fmt.Printf("wrote %d bytes\n", n)
	w.Flush()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func check_panic(e error) {
	if e != nil {
		panic(e)
	}
}
