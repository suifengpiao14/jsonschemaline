package jsonschemaline

import (
	"sort"
	"strings"
)

type NameMatch struct {
	Name       string
	Lineschema *Jsonschemaline
	Match      *Jsonschemaline
	Possible   NameMatchs
}
type NameMatchs []*NameMatch

func (ns *NameMatchs) Add(nameMatchs ...*NameMatch) {
	*ns = append(*ns, nameMatchs...)
}

func (ns NameMatchs) Names() (names []string) {
	names = make([]string, 0)
	for _, nameMatch := range ns {
		names = append(names, nameMatch.Name)
	}
	return names
}

func (ns NameMatchs) Mach(setNameMatches NameMatchs) NameMatchs {
	setNames := make([]string, 0)
	for _, setNameMatch := range setNameMatches {
		setNames = append(setNames, setNameMatch.Name)
	}
	for _, nMatch := range ns {
		setName, ok := Similarity(nMatch.Name, setNames, "")
		if ok {
			for _, setNameMatch := range setNameMatches {
				if setNameMatch.Name == setName {
					nMatch.Possible.Add(setNameMatch)
				}
			}
		}
	}
	return ns

}

func lineSchema2NameMatchs(jsonLineschema Jsonschemaline) (nameMatchs NameMatchs) {
	nameMatchs = make(NameMatchs, 0)
	baseNames := jsonLineschema.BaseNames()
	prefix := FindStrArrayCommonPrefix(baseNames, 1)
	for _, name := range baseNames {
		name = strings.TrimPrefix(name, prefix)
		nameMatch := NameMatch{
			Name:       name,
			Lineschema: &jsonLineschema,
			Possible:   make(NameMatchs, 0),
		}
		nameMatchs = append(nameMatchs, &nameMatch)
	}
	return nameMatchs

}

type StaticLineschema struct {
	Count      int
	Lineschema *Jsonschemaline
}
type StaticLineschemas []StaticLineschema

func (a StaticLineschemas) Len() int           { return len(a) }
func (a StaticLineschemas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StaticLineschemas) Less(i, j int) bool { return a[i].Count < a[j].Count }

// MatchLineschema 根据输入输出schema,从集合中匹配下游输入输出
func MatchLineschema(target Jsonschemaline, set []Jsonschemaline) (matched []Jsonschemaline) {
	targetNameMatchs := lineSchema2NameMatchs(target)
	setNameMatches := make(NameMatchs, 0)
	for _, lineschema := range set {
		nameMatches := lineSchema2NameMatchs(lineschema)
		setNameMatches.Add(nameMatches...)
	}
	nameMatches := targetNameMatchs.Mach(setNameMatches)
	_ = nameMatches
	staticsMap := make(map[*Jsonschemaline]int)
	for _, nameMatch := range nameMatches {
		for _, matched := range nameMatch.Possible {
			staticsMap[matched.Lineschema]++
		}
	}
	staticLineschemas := make(StaticLineschemas, 0)
	for lineschema, count := range staticsMap {
		sls := StaticLineschema{
			Count:      count,
			Lineschema: lineschema,
		}
		staticLineschemas = append(staticLineschemas, sls)
	}

	sort.Sort(sort.Reverse(staticLineschemas))
	for _, nameMatch := range targetNameMatchs {
		goto1 := false
		for _, sls := range staticLineschemas {
			for _, pls := range nameMatch.Possible {
				if sls.Lineschema == pls.Lineschema {
					nameMatch.Match = sls.Lineschema
					goto1 = true
					break
				}
			}
			if goto1 {
				break
			}
		}
	}
	lsm := make(map[*Jsonschemaline]struct{})
	matched = make([]Jsonschemaline, 0)
	for _, nameMatch := range targetNameMatchs {
		if _, ok := lsm[nameMatch.Match]; !ok {
			matched = append(matched, *nameMatch.Match)
			lsm[nameMatch.Match] = struct{}{}
		}
	}

	return matched
}
