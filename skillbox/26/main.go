package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func joinFilesContent(fileName1, fileName2 *string) string {
	file1Content, err := ioutil.ReadFile(*fileName1)
	errorCheck(err)

	file2Content, err := ioutil.ReadFile(*fileName2)
	errorCheck(err)

	return strings.Join([]string{string(file1Content), string(file2Content)}, "; ")
}

func printFileContent(fileName *string) {
	content, err := ioutil.ReadFile(*fileName)
	errorCheck(err)
	fmt.Printf("<%s> content: <%s>", *fileName, string(content))
}

func main() {
	switch len(os.Args) {
	case 2:
		inFileA := os.Args[1]
		printFileContent(&inFileA)
	case 3:
		inFileA := os.Args[1]
		inFileB := os.Args[2]
		content := joinFilesContent(&inFileA, &inFileB)
		fmt.Printf("<%s> + <%s> content: <%s>", inFileA, inFileB, content)
	case 4:
		inFileA := os.Args[1]
		inFileB := os.Args[2]
		outFile := os.Args[3]
		content := joinFilesContent(&inFileA, &inFileB)
		err := ioutil.WriteFile(outFile, []byte(content), 0644)
		errorCheck(err)
		fmt.Printf("Содержимое файла <%s> заполнено сконкатенированным содержанием файлов <%s> и <%s>", outFile, inFileA, inFileB)
	default:
		fmt.Println("Программу необходимо запустить с аргументами командной строки!")
		return
	}
}
