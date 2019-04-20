package main

import "fmt"

func main() {

	// Split data:
	// Open the original file url.txt that saves the URLs and read each URL by line. Then process URLs in each row.
	// Map the URL in the splitted files to < url, 1 >:
	// Read the file which saves the URL by line, and hash each line (each url) separately, then output it to some small files.
	fileRoot := "files/"
	PrepareDir(fileRoot)  // fileRoot should be empty before execute function ReadLine
	ReadLine("url.txt", processLine)
	fmt.Println("---- Finish data splitting! ----")

	// Calculate how many times the URLs have appeared in each splitted file
	mapRoot := "maps/"
	PrepareDir(mapRoot)

	err := GetMap()
	if err != nil {
		fmt.Println("生成map失败", err.Error())
	} else {
		fmt.Println("---- Finish getting map! ----")
	}

	// 读所有map文件到内存
	m, err := ReadMap()
	if err != nil {
		fmt.Println("读取map失败")
	}
	fmt.Println("---- Finish reading map! ----")

	// 将内存中的map根据值的大小进行排序
	// Sort maps (in memory) according to the size of the values 
	SortMap(m)
	fmt.Println("---- Finish sorting! ----")
}