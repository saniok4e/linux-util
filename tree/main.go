package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeIndented(out, path, printFiles, "")
}

// Вспомогательная функция с отступами
func dirTreeIndented(out io.Writer, path string, printFiles bool, indent string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	// Сортировка по имени
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	// Фильтрация: только директории или директории + файлы
	var filtered []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() || (printFiles && !entry.IsDir()) {
			filtered = append(filtered, entry)
		}
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1
		prefix := "├───"
		if isLast {
			prefix = "└───"
		}

		// Формирование строки для вывода
		name := entry.Name()
		if !entry.IsDir() {
			size := entry.Size()
			if size == 0 {
				name += " (empty)"
			} else {
				name += fmt.Sprintf(" (%db)", size)
			}
		}

		// Вывод текущего элемента
		fmt.Fprintf(out, "%s%s%s\n", indent, prefix, name)

		// Рекурсия для директорий
		if entry.IsDir() {
			nextIndent := indent + "│\t"
			if isLast {
				nextIndent = indent + "\t"
			}
			err = dirTreeIndented(out, filepath.Join(path, entry.Name()), printFiles, nextIndent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go <path> [-f]")
		os.Exit(1)
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(os.Stdout, path, printFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
