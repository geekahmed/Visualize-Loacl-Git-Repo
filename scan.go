package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

/*
 * scan given a path crawls it and its sub-folders searching for Git repositories
 */
func scan(path string){
	fmt.Printf("Found folders: \n\n")
	repos := recursiveFolderScanning(path)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repos)
	fmt.Printf("\n\nSuccessfully added\n\n")
}
func scanGitFolders(folders []string, folderPath string) []string{
	folderPath = strings.TrimSuffix(folderPath, "/")
	f, err := os.Open(folderPath)

	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil{
		log.Fatal(err)
	}
	var path string
	for _, file := range files{
		if file.IsDir(){
			path = folderPath + "/" + file.Name()
			if file.Name() == ".git"{
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules"{
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

/*
 * recursiveFolderScanning starts the recursive search of git repositories living in the `folder` subtree
*/
func recursiveFolderScanning(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

//getDotFilePath returns the dot file for the repos list.
func getDotFilePath() string{
	usr, err := user.Current()
	if err != nil{
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.gitlocalstats"
	return dotFile
}
func addNewSliceElementsToFile(filePath string, newRepos []string){
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringSliceToFile(repos, filePath)
}

func parseFileLinesToSlice(filePath string) []string{
	f := openFile(filePath)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			return nil
		}
	}
	return lines
}
func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// other error
			panic(err)
		}
	}

	return f
}

func joinSlices(newRepo []string, oldRepos []string) []string{

	for _, i := range newRepo {
		if !sliceContains(oldRepos, i){
			oldRepos = append(oldRepos, i)
		}
	}
	return oldRepos
}

func sliceContains(oldRepos []string, value string) bool{
	for _, v := range oldRepos{
		if v == value{
			return true
		}
	}
	return false
}

func dumpStringSliceToFile(repos []string, filePath string){
	content := strings.Join(repos, "\n")
	err := ioutil.WriteFile(filePath, []byte(content), 0755)
	if err != nil {
		log.Fatal(err)
	}
}