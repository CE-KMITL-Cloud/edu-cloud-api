package libvirt

import (
	etree "github.com/beevik/etree"
)

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
