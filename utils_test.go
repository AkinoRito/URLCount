package main

import "testing"
import "sort"
// import "strconv"
import "os"

func TestReadLine(t *testing.T) {
	var filePath = "url.txt"
	err := ReadLine(filePath, processLine)
	if err != nil {
		t.Error(err)
	}
	
}


// Test functions related with MapSorter 
func TestSortMap(t *testing.T) {
	ms := make(MapSorter, 0, 5)
	ms = append(ms, Item{"0", 0})
	ms = append(ms, Item{"3", 3})
	ms = append(ms, Item{"2", 2})
	ms = append(ms, Item{"1", 1})
	ms = append(ms, Item{"4", 4})
	sort.Sort(ms)

	for i, item := range ms {
		if i + item.Val != 4 {
			t.Errorf("Wrong sorting result!")
		}
	}
}


func TestPrepareDir(t *testing.T) {
	// 1. Use an empty folder as a test.
	emptyPath := "test/emptyPath"
	_, err := os.Stat(emptyPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(emptyPath, os.ModePerm)
		} else {
			t.Errorf("ERROR: Empty folder create error")
		}
	}

	err = PrepareDir(emptyPath)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = os.Stat(emptyPath)
	if err != nil {
		t.Errorf("ERROR: PrepareDir error when the folder is empty")
	}


	// 2. Use a nonexistant folder as a test.
	nonExistPath := "test/nonExistPath"
	os.RemoveAll(nonExistPath)

	err = PrepareDir(nonExistPath)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = os.Stat(nonExistPath)
	if err != nil {
		t.Errorf("ERROR: PrepareDir error when the folder is nonexistant")
	}


	// 3. Use a non-empty folder as a test
	nonEmptyFolder := "test/nonEmptyFolder"
	_, err = os.Stat(nonEmptyFolder)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(nonEmptyFolder, os.ModePerm)
			f, errFileCreate := os.Create("newFile")
			defer f.Close()
			if errFileCreate != nil {
				t.Errorf(errFileCreate.Error())
			}
		} else {
			t.Errorf("ERROR: Empty folder create error")
		}
	}

	err = PrepareDir(nonEmptyFolder)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = os.Stat(nonEmptyFolder)
	if err != nil {
		t.Errorf("ERROR: PrepareDir error when the folder is non-empty")
	}
}