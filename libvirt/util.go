package libvirt

import (
	"bufio"
	"log"
	"os"

	etree "github.com/beevik/etree"
)

type UEFIArch struct {
	i686    []string
	x86_64  []string
	aarch64 []string
	armv7l  []string
}

func GetXPath(file string, path string) (string, error) {
	doc := etree.NewDocument()
	var result string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for _, e := range doc.FindElements(path) {
		result = e.Text()
	}
	return result, nil
}

func WriteStringtoFile(input string, output_file string) {
	f, err := os.Create(output_file)
	check_panic(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(input)
	// check_panic(err)
	// fmt.Printf("wrote %d bytes\n", n)
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

func UEFIArchPatterns() UEFIArch {
	return UEFIArch{
		i686: []string{
			`.*ovmf-ia32.*`, // fedora, gerd's firmware repo
		},
		x86_64: []string{
			`.*OVMF_CODE\.fd`,       // RHEL
			`.*ovmf-x64/OVMF.*\.fd`, // gerd's firmware repo
			`*ovmf-x86_64-.*`,       // SUSE
			`.*ovmf.*`,
			`.*OVMF.*`, // generic attempt at a catchall
		},
		aarch64: []string{
			`.*AAVMF_CODE\.fd`,     // RHEL
			`.*aarch64/QEMU_EFI.*`, // gerd's firmware repo
			`.*aarch64.*`,          // generic attempt at a catchall
		},
		armv7l: []string{
			`.*arm/QEMU_EFI.*`, // fedora, gerd's firmware repo
		},
	}
}
