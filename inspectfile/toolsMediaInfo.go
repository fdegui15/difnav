// ToolMediainfo for "inspectFile"
// use Mediainfo
// https://mediaarea.net/fr/MediaInfo

package main

import (
	"encoding/json"
	"fmt"
)

func init() {
	//Integration of the new tool
	tools["mi"] = toolsStr{
		Cmd:  "mediainfo",
		Name: "mediainfo",
		//Args: []string{"--Inform=file:///home/go/src/inspectfile/MediaInfoTemplate.txt"},
		Args: []string{},
		Vers: []string{"--version"},
		fn:   launchMediaInfo}
}

func launchMediaInfo(fbyte []byte) []byte {
	toolFlag := "mi"

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	sfstring, _ := exectoolsCmd(toolFlag, fn)

	fmt.Println("mediainfo: Output=" + sfstring)

	sfm := make(map[string]interface{}) // Category
	sfmint := make(map[string]string)   // Element in the Category

	arrstrg := strings.Split(sfstring, "\r") // Split on new line
	var catstrg string                       // Category Name

	for i := 0; i < len(arrstrg); i++ {
		line := strings.Trim(arrstrg[i], " \r\n")
		if len(line) > 0 {
			if !(strings.Contains(line, ":")) {
				//New Category
				if len(sfmint) > 0 {
					//if sfmint is not empty, let's save the old Category
					//fmt.Println("sfm["+catstrg+"]= ", sfmint, len(sfmint))
					sfm[catstrg] = sfmint
				}
				catstrg = line
				sfmint = make(map[string]string) // Initialise sfmint
			} else {
				// Save the new element in the Category
				varstrg := strings.SplitN(line, ":", 2)
				//fmt.Println("Val[" + catstrg + "][" + strings.Trim(varstrg[0], " \n") + "] = \"" + strings.Trim(varstrg[1], " \n") + "\"")
				sfmint[strings.Trim(varstrg[0], " ")] = strings.Trim(varstrg[1], " ")
			}
		}
	}
	if len(sfmint) > 0 {
		//Last Category !
		sfm[catstrg] = sfmint
	}

	//Define this ouput here
	f[tools[toolFlag].Name] = sfm

	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}
