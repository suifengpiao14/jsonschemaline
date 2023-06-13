package jsonschemaline

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/suifengpiao14/funcs"
)

// jsonschemaline 生成go 结构体工具
type Struct struct {
	IsRoot     bool
	Name       string
	Lineschema string
	Attrs      []*StructAttr
}

// AddAttrIgnore 已经存在则跳过
func (s *Struct) AddAttrIgnore(attrs ...StructAttr) {
	if len(s.Attrs) == 0 {
		s.Attrs = make([]*StructAttr, 0)
	}
	for _, attr := range attrs {
		if _, exists := s.GetAttr(attr.Name); exists {
			continue
		}
		s.Attrs = append(s.Attrs, &attr)
	}
}

// AddAttrReplace 增加或者替换
func (s *Struct) AddAttrReplace(attrs ...StructAttr) {
	if len(s.Attrs) == 0 {
		s.Attrs = make([]*StructAttr, 0)
	}
	for _, attr := range attrs {
		if old, exists := s.GetAttr(attr.Name); exists {
			*old = attr
			continue
		}
		s.Attrs = append(s.Attrs, &attr)
	}
}
func (s *Struct) GetAttr(attrName string) (structAttr *StructAttr, exists bool) {
	for _, attr := range s.Attrs {
		if attr.Name == attrName {
			return attr, true
		}
	}
	return nil, false
}

type StructAttr struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type Structs []*Struct

func (s *Structs) Json() (str string) {
	b, _ := json.Marshal(s)
	str = string(b)
	return str
}

func (s *Structs) GetRoot() (struc *Struct, exists bool) {
	for _, stru := range *s {
		if stru.IsRoot {
			return stru, true
		}
	}
	return struc, false
}

func (s *Structs) Get(name string) (struc *Struct, exists bool) {
	for _, stru := range *s {
		if stru.Name == name {
			return stru, true
		}
	}
	return struc, false
}

func (s *Structs) AddIngore(structs ...*Struct) {
	if len(*s) == 0 {
		*s = make(Structs, 0)
	}
	for _, structRef := range structs {
		if _, exists := s.Get(structRef.Name); exists {
			continue
		}
		*s = append(*s, structRef)

	}
}

// Copy 深度复制
func (s Structs) Copy() (newStructs Structs) {
	newStructs = make(Structs, 0)
	for _, struc := range s {
		newStruct := *struc
		newStruct.Attrs = make([]*StructAttr, 0)
		for _, attr := range struc.Attrs {
			newAttr := *attr
			newStruct.Attrs = append(newStruct.Attrs, &newAttr)
		}
		newStructs = append(newStructs, &newStruct)
	}
	return newStructs
}

func (s *Structs) AddNameprefix(nameprefix string) {
	if len(*s) == 0 {
		return
	}
	allAttrs := make([]*StructAttr, 0)
	for _, struc := range *s {
		allAttrs = append(allAttrs, struc.Attrs...)
	}
	for _, struc := range *s {
		baseName := struc.Name
		struc.Name = funcs.ToCamel(fmt.Sprintf("%s_%s", nameprefix, baseName))
		for _, attr := range allAttrs {
			if strings.HasSuffix(attr.Type, baseName) {
				attr.Type = fmt.Sprintf("%s%s", attr.Type[:len(attr.Type)-len(baseName)], struc.Name)
			}
		}
	}
}
