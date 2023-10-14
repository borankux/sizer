package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type folderInfo struct {
	name string
	size int64
}

type bySizeDesc []folderInfo

func (a bySizeDesc) Len() int           { return len(a) }
func (a bySizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySizeDesc) Less(i, j int) bool { return a[i].size > a[j].size }

func main() {
	folderPath := "." // Change this to the path of the folder you want to scan

	folders, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var folderList []folderInfo

	for _, folder := range folders {
		if folder.IsDir() {
			folderSize, err := calculateFolderSize(filepath.Join(folderPath, folder.Name()))
			if err != nil {
				fmt.Printf("Error calculating size of %s: %v\n", folder.Name(), err)
				continue
			}
			folderList = append(folderList, folderInfo{name: folder.Name(), size: folderSize})
		}
	}

	sort.Sort(bySizeDesc(folderList))

	for _, folder := range folderList {
		// Format the size with a fixed width of 12 characters
		sizeStr := fmt.Sprintf("%12s", formatSize(folder.size))
		fmt.Printf("%s %s\n", sizeStr, folder.name)
	}
}

func calculateFolderSize(folderPath string) (int64, error) {
	var folderSize int64

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			folderSize += info.Size()
		}
		return nil
	})

	return folderSize, err
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}
