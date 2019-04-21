package main

import "fmt"

func main() {

	// Split data:
	// Open the original file url.txt that saves the URLs and read each URL by line. Then process URLs in each row.
	// Map the URL in the splitted files to < url, 1 >:
	// Read the file which saves the URL by line, and hash each line (each url) separately, then output it to some small files.
	// fileRoot := "files/"
	// PrepareDir(fileRoot)  // fileRoot should be empty before execute function ReadLine
	// ReadLine("url.txt", processLine)
	fmt.Println("---- Finish data splitting! ----")

	// Calculate how many times the URLs have appeared in each splitted file
	mapRoot := "maps/"
	PrepareDir(mapRoot)

	err := GetMap()
	if err != nil {
		fmt.Println("Failed to generate map", err.Error())
	} else {
		fmt.Println("---- Finish getting map! ----")
	}

	// Read all map files to memory
	m, err := ReadMap()
	if err != nil {
		fmt.Println("Failed to read map")
	}
	fmt.Println("---- Finish reading map! ----")

	// Sort maps (in memory) according to the size of the values 
	SortMap(m)
	fmt.Println("---- Finish sorting! ----")
}