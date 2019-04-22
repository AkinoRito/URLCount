package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"crypto/md5"
	"encoding/binary"
	"strconv"
	"bytes"
	"strings"
	"sort"
)

var ChunkNum uint64 = 150
var MapFileNum int = 101


func processLine(line []byte) {
	/*
	* Divide the raw file into 150 small files.
	* Rules: 
	* Calculate the hash value of each URL and divide the hash value by 150 to get the file name of the output files
	* In this case, each URL is split into the same file . 
	*/
	hash := md5.New()
	_, err := hash.Write(line)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	result := hash.Sum(nil)  // type of `result`: uint8
	data := binary.BigEndian.Uint64(result)
	dataInt := int(data % ChunkNum)
	fileName := "files/file_" + strconv.Itoa(dataInt) + ".txt"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return //err
	}
	defer f.Close()

	// If a URL is too short, it may not be a correct URL
	lineLen := len(line)
	if lineLen < 5 {
		return
	}

	// If the file is small, you can write the URL to the corresponding file after a hash; 
	// if the file corresponding to one hash is too large, which will exceed 1GB after writing, then the URL will be hashed again.
	for {
		fileLen, _ := f.Seek(0, os.SEEK_END)
		if lineLen + int(fileLen) < 2 << 30  {
			break
		} else {
			dataInt = (dataInt + 1) % MapFileNum
			fileName = "files/file_" + strconv.Itoa(dataInt) + ".txt"
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return //err
			}
			defer f.Close()

			fileLen, _ = f.Seek(0, os.SEEK_END)
		}
	}
	f.Write(line)

}


func ReadLine(filePath string, processLine func([]byte)) error {
	/*
	* Read each line of a file
	*/
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		processLine(line)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
	return nil
}


func GetMap() error {
	/*
	* Calculate how many times the URLs have appeared in each splitted file.
	* And the result will be recorded in `map_xx.txt`
	*/

	for i := 0; i < int(ChunkNum); i++ {
		// Open a splitted URL file
		fileName := "files/file_" + strconv.Itoa(int(i)) + ".txt"
		f, err := os.Open(fileName)

		if err != nil {
			// fmt.Println(fileName + " does not exist")
			continue
			// return err
		}
		defer f.Close()

		m := make(map[string]int)
		// Traverse all URLs in each file and use map to count the number of times each URL appears 
		bfRd := bufio.NewReader(f)
		for {
			line, err := bfRd.ReadBytes('\n')
			line = bytes.TrimRight(line, "\r\n")
			// Update the map result
			lineStr := string(line[:])
			if lineStr == "" {
				break
			} else {
				if m[lineStr] != 0 {
					m[lineStr]++
				} else {
					m[lineStr] = 1
				}
			}
			
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}

		// Save map into a file
		mapNum := int(i) % MapFileNum
		mapFileName := "maps/map_" + strconv.Itoa(mapNum) + ".txt"
		mapF, errMapF := os.OpenFile(mapFileName, os.O_APPEND|os.O_CREATE, 0644)
		if errMapF != nil {
			// fmt.Println("Failed to open " + mapFileName)
			continue
		}
		defer mapF.Close()

		for k, v := range m {
			// Calculate the number of bits to be written to the map.
			mapBytes := []byte(k + "," + strconv.Itoa(v) + "\r\n")
			mapLen := len(mapBytes)
			
			// Calculate the number of bytes already in the file 
			// Calculate the number of bits in an existing map file
			fileLen, _ := mapF.Seek(0, os.SEEK_END)

			if mapLen + int(fileLen) < 2 << 30 {
				mapF.Write(mapBytes)
			} else {
				for {
					if mapLen + int(fileLen) < 2 << 30 {
						break
					} else {
						// Re-hash the key-value pairs and save them to another map file.
						mapNum = (mapNum + 1) % MapFileNum
						mapFileName = "maps/map_" + strconv.Itoa(mapNum) + ".txt"
						mapF, errMapF := os.OpenFile(mapFileName, os.O_APPEND|os.O_CREATE, 0644)
						defer mapF.Close()
						if errMapF != nil {
							fmt.Println("Failed to open " + mapFileName)
							continue
						}
						fileLen, _ = mapF.Seek(0, os.SEEK_END)
					}
				}
				
			}
			
		}

	}
	return nil
}


func ReadMap() (map[string]int, error){
	// Read all map files to memory
	m := make(map[string]int)

	for i := 0; i < MapFileNum; i++ {
		mapFileName := "maps/map_" + strconv.Itoa(int(i) % MapFileNum) + ".txt"
		f, err := os.Open(mapFileName)
		if err != nil {
			// return err
			continue
		}
		defer f.Close()

		bfRd := bufio.NewReader(f)
		for {
			line, err := bfRd.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
			mapLine := string(bytes.TrimRight(line, "\r\n"))
			list := strings.Split(mapLine, ",")
			if m[list[0]] != 0 {
				return nil, err
			} else {
				m[list[0]], _ = strconv.Atoi(list[1])
			}
		}
	}
	return m, nil
}


/* Build a map sorter */
type Item struct {
	Key string
	Val int
}
type MapSorter []Item


func NewMapSorter(m map[string]int) MapSorter {
	ms := make(MapSorter, 0, len(m))
	for k, v := range m {
		ms = append(ms, Item{k, v})
	}
	return ms
}

// Implement sort.Interface (implement Swap method, Len method end Less method)
func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	// Rank the map by value
	return ms[i].Val > ms[j].Val
}


func SortMap(m map[string]int) {
	/*
	* Sort maps (in memory) according to the size of the values in descending order/c
	*/
	ms := NewMapSorter(m)
	sort.Sort(ms)
	for i, item := range ms {
		if i >= 100 {
			break
		}
		fmt.Println("top", i+1, ":", item.Key, item.Val)
	}
}


func PrepareDir(path string) error {
	/*
	* If the directory does not exist, create an empty directory; 
	* if the directory exists but is not an empty directory, empty the directory.
	*/

	// Determine if the folder exists
	_, err := os.Stat(path)
	if err == nil {
		// The folder exists, then it should be ensured that the folder is empty
		dir, _ := ioutil.ReadDir(path)
		if len(dir) == 0 {
			fmt.Println(path, "is empty")
			return nil
		} else {
			fmt.Println(path, "is not empty, it will be cleared now.")
			os.RemoveAll(path)  // Function os.RemoveAll returns error only when the path is not exist, but this situation won't appear here.
			os.Mkdir(path, os.ModePerm)
		}
	} else if os.IsNotExist(err) {
		// Create a dnew folder
		fmt.Println(path, "does not exist, it will be created now.")
		os.Mkdir(path, os.ModePerm)
	} else {
		return err
	}
	return nil
}
