package internal

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	ID           string   `json:"id"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
	Folder       string   `json:"folder"`
	Extension    string   `json:"extension"`
	Name         string   `json:"name"`
	Lines        int      `json:"lines"`
	Label        string   `json:"label"`
}
type Label struct {
	Pattern *regexp.Regexp
	Name    string
}

type Config struct {
	Root              string
	FilesPattern      *regexp.Regexp
	Labels            []Label
	labelByFolderName bool
}
type InterfaceParser interface {
	getMatchingFilePaths() ([]string, error)
	getFileDependencies(filePath string, deps chan []string) ([]string, error)
	init()
}
type Parser struct {
	Config
}

type StringTransformer func(string) string

func NewConfig(root string, pattern string, labels []Label) Parser {
	return Parser{
		Config{
			Root:         root,
			FilesPattern: regexp.MustCompile(pattern),
			//TODO:  will label specific patterns
			Labels: labels,
			labelByFolderName: true,
		}}
}

func (parser *Parser) getMatchingFilePaths() ([]string, error) {
	var (
		root  = parser.Root
		regex = parser.FilesPattern
	)
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isSupportedPath := regex.MatchString(path); !isSupportedPath {
			return nil
		}
		if info.IsDir() {
			return nil
		} else {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (parser *Parser) getFileDependencies(filePath string, deps chan []string) ([]string, error) {
	importRegex := `import(?:["'\s]*([\w*${}\n\r\t, ]+)from\s*)?["'\s]["'\s](.*[@\w_-]+)["'\s].*;?$`

	f, e := os.Open(filePath)
	defer f.Close()

	if e != nil {
		panic(e)
	}
	b := new(strings.Builder)
	io.Copy(b, f)
	r, _ := regexp.Compile("(?m)" + importRegex)

	res := r.FindAllString(b.String(), -1)
	importPathReges, _ := regexp.Compile(`['"].*['"]`)
	dependencies := make([]string, 0)
	for _, s := range res {
		path := importPathReges.FindString(s)
		dependencies = append(dependencies, path)
	}

	deps <- dependencies
	return dependencies, nil
}


func getFileNameWithoutExtension(path string, extension string) string {
	fileNameWithExtensionRegex, _ := regexp.Compile(`(?m)[A-Za-z0-9_\-\.]+\.[A-Za-z0-9]+$`)
	name := fileNameWithExtensionRegex.FindString(path)
	return strings.ReplaceAll(name, "."+extension, "")
	//return strings.Trim(name, "."+extension)

}
func getFileExtension(fileName string) string {
	extensionRegex, _ := regexp.Compile(`\w+$`)
	return extensionRegex.FindString(fileName)
}

func removeQuotesFromPath(path string) string {
	reg, _ := regexp.Compile(`(\w?(-)|\w|\.|\.\.|\/)+(\w*)`)
	return reg.FindString(path)
}

func transformPaths(paths []string, transformer StringTransformer) []string {
	var transformedPaths = make([]string, 0)
	relativeReg, _ := regexp.Compile(`^\.`)
	for _, path := range paths {
		path = transformer(path)

		isRelative := relativeReg.MatchString(path)

		if isRelative {
			reg, _ := regexp.Compile(`[A-Za-z0-9_\-\.]+$`)
			path = reg.FindString(path)
		}
		transformedPaths = append(transformedPaths, path)
	}
	return transformedPaths
}

func (parser *Parser) Init(relativePath string) ([]File, map[string]File, map[string][]string) {
	fileMap := make(map[string]File)
	dependenciesMap := make(map[string][]string)

	files, err := parser.getMatchingFilePaths()

	if err != nil {
		log.Fatal(err)
	}

	jsonArr := make([]File, 0)

	for _, file := range files {
		deps := make(chan []string)
		re := regexp.MustCompile(`.+\/(.+)\/[^\/]*`)
		arr := re.FindStringSubmatch(file)
		openedFile, _ := os.Open(file)
		defer openedFile.Close()
		lines, _ := lineCounter(bufio.NewReader(openedFile))

		go parser.getFileDependencies(file, deps)
		value := <-deps
		fileInfo := File{
			ID:           GetUniqueID(),
			Lines:        lines,
			Path:         file,
			Dependencies: transformPaths(value, removeQuotesFromPath),
			Extension:    getFileExtension(file),
			Folder:       relativePath,
			Name:         getFileNameWithoutExtension(file, getFileExtension(file)),
		}

		if len(arr) == 2 {
			fileInfo.Folder = arr[1]
		}

		fileMap[fileInfo.Name] = fileInfo
		dependenciesMap[fileInfo.Name] = fileInfo.Dependencies
		jsonArr = append(jsonArr, fileInfo)

	}
	return jsonArr, fileMap, dependenciesMap

}
