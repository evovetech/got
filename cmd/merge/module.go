package merge

import (
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/log"
	"regexp"
)

type ModuleFile struct {
	ProjectFile

	Module DirPath
}

func parseSrc(fp FilePath) *ModuleFile {
	if match := reSrc.FindStringSubmatch(fp.slashy); match != nil {
		src := new(ModuleFile)
		src.ProjectFile = *newProjectFile(fp, match[1])
		src.Module = src.RelPath.ToDir()
		src.RelPath = GetFilePath(match[2])
		return src
	}
	return nil
}

type Module struct {
	Project DirPath      `json:"-"`
	Name    DirPath      `json:"-"`
	Re      string       `json:",omitempty"`
	Src     []ModuleFile `json:"-"`
	Other   []ModuleFile `json:",omitempty"`

	re *regexp.Regexp
}

func (m *Module) parse(pf ProjectFile) *ModuleFile {
	if m.re == nil {
		name := m.Name.slashy
		if name == "." {
			name = ""
		} else {
			name += "/?"
		}
		pat := fmt.Sprintf("^%s(.*)$", name)
		m.re = regexp.MustCompile(pat)
		m.Re = m.re.String()
	}
	fp := pf.RelPath
	if match := m.re.FindStringSubmatch(fp.slashy); match != nil {
		f := new(ModuleFile)
		f.ProjectFile = pf
		f.Module = m.Name
		f.RelPath = GetFilePath(match[1])
		return f
	}
	log.Printf("\\%s\\unable to parse '%s'", m.re, fp)
	return nil
}

func (m *Module) addSrc(file ModuleFile) {
	m.Src = append(m.Src, file)
}

func (m *Module) addOther(file ModuleFile) {
	m.Other = append(m.Other, file)
}

type Modules map[DirPath]*Module

func (mods Modules) MarshalJSON() ([]byte, error) {
	modules := make(map[string]*Module)
	for k, v := range mods {
		modules[k.slashy] = v
	}
	return json.Marshal(modules)
}
