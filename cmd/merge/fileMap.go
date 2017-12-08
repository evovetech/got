package merge

import "fmt"

type FileMap struct {
	Name     string
	Path     FilePath
	Children map[string]*FileMap
}

var RootFile = GetFileMap("")

func GetFileMap(file string) *FileMap {
	fp := GetFilePath(file)
	return &FileMap{
		Name:     fp.Name(),
		Path:     fp,
		Children: make(map[string]*FileMap),
	}
}

func (m *FileMap) IsDir() bool {
	return m.Path.IsDir()
}

func (m *FileMap) Add(file FilePath) error {
	if m.Path.Name() != "." {
		return fmt.Errorf("can only add to root, not %s", m)
	}
	var dir DirPath
	if file.IsDir() {
		dir = GetDirPath(file.actual)
	} else {
		dir = file.Dir()
	}
	if dir.Name() == "." {
		m.Children[file.Name()] = GetFileMap(file.actual)
	} else {
		// TODO:
	}
	return nil
}

func (m *FileMap) String() string {
	return m.Path.String()
}
