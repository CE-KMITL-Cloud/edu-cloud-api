package libvirt

import (
	"bufio"
	"log"
	"os"

	etree "github.com/beevik/etree"
)

func GetXPath(file string, path string) (string, error) {
	doc := etree.NewDocument()
	var result string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	e := doc.FindElement(path)
	if e != nil {
		result = e.Text()
	} else {
		result = ""
	}
	return result, nil
}

func GetElementsLength(file string, path string) int {
	doc := etree.NewDocument()
	var arr []string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for _, e := range doc.FindElements(path) {
		arr = append(arr, e.Text())
	}
	return len(arr)
}

func GetChildElementsAttr(file string, path string, value []string, key string) ([]string, error) {
	doc := etree.NewDocument()
	var result []string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for i, e := range doc.FindElements(path) {
		if e.Text() == value[i] {
			if e.SelectAttr(key) != nil {
				result = append(result, e.SelectAttr(key).Value)
			} else {
				result = append(result, "")
			}
		}
	}
	return result, nil
}

func GetXPathsAttr(file string, path string, key string) ([]string, error) {
	doc := etree.NewDocument()
	length := GetElementsLength(file, path)
	result := make([]string, length)
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for i, e := range doc.FindElements(path) {
		result[i] = e.SelectAttr(key).Value

	}
	return result, nil
}

func GetXPaths(file string, path string) ([]string, error) {
	doc := etree.NewDocument()
	length := GetElementsLength(file, path)
	result := make([]string, length)
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for i, e := range doc.FindElements(path) {
		result[i] = e.Text()
	}
	return result, nil
}

func GetParentTags(file string, path string, element string) ([]string, error) {
	doc := etree.NewDocument()
	var result []string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	for _, e := range doc.FindElements(path) {
		for j := range e.Parent().SelectElement(element).ChildElements() {
			result = append(result, e.Parent().SelectElement(element).ChildElements()[j].Tag)
		}
	}
	return result, nil
}

func GetParentText(file string, path string, element string) (string, error) {
	doc := etree.NewDocument()
	var result string
	if err := doc.ReadFromFile(file); err != nil {
		log.Fatalln(err)
	}
	e := doc.FindElement(path)
	if e != nil {
		result = e.Parent().SelectElement(element).Text()
	} else {
		result = ""
	}
	return result, nil
}

func WriteStringtoFile(input string, output_file string) {
	f, err := os.Create(output_file)
	check_panic(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(input)
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
