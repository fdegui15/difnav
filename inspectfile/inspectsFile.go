// Main module for "inspectFile"

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//Version number of InspectFile
var Version = [3]int{0, 1, 0}

//output of InspectFile (in json format)
var output []byte

//flags
var (
	dir     = flag.Bool("dir", false, "inspects the directory, sub included.")
	version = flag.Bool("version", false, "display version information")
	flserve = flag.String("server", "", "start inspects server e.g. -server localhost:5138")
	fltools = flag.String("tools", "", "includes the tools you want to use, if more than 2, it must be comma separated (hash,fido,sf,et,mi)")
)

func inspectfile(id int, filename string, fbyte []byte) []byte {
	type FileStr struct {
		_ProcessId    int
		FileName      string
		FileShortName string `json:",omitempty"`
		FileSize      int64  `json:",omitempty"`
		FileExists    bool
		FileNotNull   bool
		FileComments  string `json:",omitempty"`
	}
	f := FileStr{
		FileName:    filename,
		FileSize:    -1,
		FileExists:  false,
		FileNotNull: false,
		_ProcessId:  id,
	}
	shortnTab := strings.Split(f.FileName, "/")
	f.FileShortName = shortnTab[len(shortnTab)-1]

	var fcurbyte []byte //, _ := json.MarshalIndent(f, "", "    ") //To initialise the []byte
	finfo, err := os.Stat(filename)
	if err != nil {
		f.FileComments = "the file doesn't exist..."
		fcurbyte, _ = json.MarshalIndent(f, "", "    ")
		fmt.Printf(msgWithDateProc(id, "inspectfile - the file doesn't exists:  "+filename))
	} else if finfo.IsDir() {
		f.FileComments = "this is a directory!!!"
		fcurbyte, _ = json.MarshalIndent(f, "", "    ")
		fmt.Println(msgWithDateProc(id, "inspectfile - this is a directore: "+filename))
	} else {
		// it's a file
		f.FileSize = finfo.Size()
		f.FileExists = true
		f.FileNotNull = true

		//str, _ = antiVirus(os.Args[1])
		fcurbyte, _ = json.MarshalIndent(f, "", "    ")

		//execute the tool in the listtools
		for _, toolflag := range listtools {
			fcurbyte = tools[toolflag].fn(id, fcurbyte)
		}
	}
	return append(fbyte, fcurbyte...)
}

func walkdir(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	if f.IsDir() == false {
		output = inspectfile(-1, path, output) //How to pass ProcessId through filepath.Walk ???
	}
	return nil
}

func inspectdir(id int, path string) []byte {
	err := filepath.Walk(path, walkdir)
	if err != nil {
		fmt.Printf(msgWithDateProc(id, "Error in parsing directory: "+err.Error()))
		return nil
	}
	return output
}

func initTools(id int, sttools string) {
	//Init listtools
	var outst string
	if sttools != "" {
		//use only the tools in the command line
		listtools = strings.Split(sttools, "-")
		outst = sttools
		//TO DO: improve to test that the tool flag exists !!!
	} else {
		//use all the tools
		i := 0
		for key, _ := range tools {
			listtools[i] = key
			i += 1
			outst += key + "-"
		}
	}
	fmt.Println(msgWithDateProc(id, "inittools: "+outst))
}

func main() {
	start := time.Now()

	flag.Parse()

	if *version {
		fmt.Printf("inspectFile %d.%d.%d\n", Version[0], Version[1], Version[2])
		return
	}
	if *flserve != "" {
		fmt.Printf("Startinf server at %s. Use CTRL+C to quit.\n", *flserve)
		listen(*flserve)
		return
	}
	if *fltools != "" {
		initTools(0, *fltools)
	}
	if *dir {
		if flag.NArg() != 1 {
			fmt.Println("You must pass a directory name!!!")
			return
		}
		inspectdir(0, flag.Arg(0))
	} else if flag.NArg() != 1 {
		fmt.Println("You must pass a file name in parameter!!!")
		return
	} else {
		output = inspectfile(0, flag.Arg(0), nil)
	}
	//output = inspectsfile("/media/sf_Temp/Benchmark.pptx", output)
	fmt.Printf("Output: %s \n", output)
	fmt.Printf("ToolsVersion: %s\n", ExportToolsVersion())

	fmt.Printf("Took %v to run.\n", time.Since(start))
}
