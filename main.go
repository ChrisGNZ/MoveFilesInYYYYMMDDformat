package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    //"path/filepath"
    "regexp"
    //"strings"
    //"time"
    )

func main() {

	// the folder to examine should be passed as a command line argument or parameter.  If no parameters supplied, then default to current folder
	folderPath := getAndCheckFolderPath()

	// obtain the contents of the folderpath 
	fmt.Println("Scanning folder ", folderPath, " for files in *YYYYMMDD*.* format...")
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		// an unexpected error occurred while obtaining the list of files in the folder path...
		log.Fatal(err)
	}


	// Compile the regex search expression once in advance to improve performance
	// Use raw strings to avoid having to quote the backslashes.
	var YYYYMMDDexp = regexp.MustCompile(`([12]\d{3}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01]))`)


	movedFileCount := 0

	//now scan the files for files that match the YYYYMMDD pattern
	// and check it is a file and not a folder
	// and if ok then create the subfolder if required and move the file to the subfolder
	for _, f := range files {
		fileName := f.Name()
		if fileNameIsInYYYYMMDD_Format( fileName, YYYYMMDDexp ) { 

			if fileIsNotFolder( folderPath+"/"+fileName ) {

				moveFileToYYYYMMDD_Subfolder( folderPath, fileName, YYYYMMDDexp )
				movedFileCount++ 
			}
		}
	}

	fmt.Println("Scan complete! Found and moved ", movedFileCount, " files in *YYYYMMDD*.* format.")
}


/////////////////////////////////////////////////////////////////////////////////////////////////////
func fileNameIsInYYYYMMDD_Format(filename string, precompiledExpression *regexp.Regexp) bool {

	return precompiledExpression.MatchString(filename) 
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
func fileIsNotFolder(filespec string ) bool {
	src, err := os.Stat(filespec)
	if err != nil {
		fmt.Println(filespec, " not found! Quiting...")
		os.Exit(1)
	}

	return !src.IsDir()  //return true if file is not a folder else return false
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
func moveFileToYYYYMMDD_Subfolder(folderPath string, fileName string, precompiledExpression *regexp.Regexp){

	// first, build the name of the subfolder in YYYY_MM_DD format 
	loc := precompiledExpression.FindIndex([]byte(fileName))
	matchedname := string(fileName[loc[0]:loc[1]])

	if len(matchedname) < 8 {
		return
	}

	YYYY := matchedname[0:4]
	MM := matchedname[4:6]
	DD := matchedname[6:8]
	subfoldername := YYYY+"_"+MM+"_"+DD


	// second, check if a folder for this timestamp has already been created?
	err := CreateDirIfNotExist(folderPath + "/" + subfoldername)
	if err != nil {
		fmt.Println("Error! Unable to create sub-folder! Giving Up! The Error Message is: ", err)
		os.Exit(2)
	}

	oldpath := folderPath + "/" + fileName
	newpath := folderPath + "/" + subfoldername + "/" + fileName

	err = os.Rename(oldpath, newpath)
	if err != nil {
		fmt.Println("Error! ", err)
		os.Exit(3)
	}
	fmt.Println(fileName, "moved to", subfoldername)

}
/////////////////////////////////////////////////////////////////////////////////////////////////////
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
func getAndCheckFolderPath() string {
	argsWithoutProg := os.Args[1:]

	// default to current folder unless supplied as a command line argument
	folderpath := "./"
	if len(argsWithoutProg) != 0 {
		folderpath = argsWithoutProg[0]
	}

	//check if the folder exists
	src, err := os.Stat(folderpath)
	if err != nil {
		fmt.Println("Folder: ", folderpath, " not found")
		os.Exit(1)
	}

	//check if the folderpath is indeed a folder and not a file
	if !src.IsDir() {
		fmt.Println(folderpath, " is not a folder")
		os.Exit(1)
	}
	return folderpath
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
