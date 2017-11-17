package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jeffail/gabs"
)

func main() {

	templatePath := flag.String("template", "", "Path to the template JSON file")
	comparePath := flag.String("file", "", "Path to the JSON file to compare to")
	flag.Parse()

	//check if template path exists
	if _, err := os.Stat(*templatePath); os.IsNotExist(err) {
		fmt.Printf("Template file path invalid.  No file found at %v\n", *templatePath)
		os.Exit(0)
	}

	if _, err := os.Stat(*comparePath); os.IsNotExist(err) {
		fmt.Printf("No file found to compare at %v\n", *comparePath)
	}

	templateFile, e := ioutil.ReadFile(*templatePath)
	if e != nil {
		fmt.Printf("Error occurred while opening template file %v\n", e)
		os.Exit(0)
	}

	compareFile, e := ioutil.ReadFile(*comparePath)
	if e != nil {
		fmt.Printf("Error occurred while opening template file %v\n", e)
		os.Exit(0)
	}

	templateParsed, err := gabs.ParseJSON(templateFile)
	if err != nil {
		fmt.Printf("Error occurred while parsing template file %v\n", err)
		os.Exit(0)
	}
	templateChildren, err := templateParsed.ChildrenMap()
	if err != nil {
		fmt.Printf("Error occurred while parsing template file %v\n", err)
		os.Exit(0)
	}
	templateList := flattenJSON(templateChildren, "")

	compareParsed, err := gabs.ParseJSON(compareFile)
	if err != nil {
		fmt.Printf("Error occurred while parsing compare file %v\n", err)
		os.Exit(0)
	}
	compareChildren, err := compareParsed.ChildrenMap()
	if err != nil {
		fmt.Printf("Error occurred while parsing compare file %v\n", err)
		os.Exit(0)
	}
	compareList := flattenJSON(compareChildren, "")

	for key, tempVal := range templateList {
		if compVal, ok := compareList[key]; ok {
			//switch on the template values
			switch tvCast := tempVal.(type) {
			case []interface{}:
				//if we have an array interface we need to compare the values so cast the compare as well
				switch cvCast := compVal.(type) {
				case []interface{}:
					//do comparison of interface arrays
					for _, tvCastVal := range tvCast {
						foundVal := false
						for _, cvCastVal := range cvCast {
							if cvCastVal == tvCastVal {
								foundVal = true
							}
						}
						//if we dont find a match, print the missing value
						if !foundVal {
							fmt.Printf("%v does not match! Template contains %v and compare is missing\n", key, tvCastVal)
						}
					}
				default:
					fmt.Printf("%v does not match! Template has %v as array and compare has %v\n", key, tempVal, compVal)
				}
			default:
				if compVal != tempVal {
					fmt.Printf("%v does not match! Template has %v and compare has %v\n", key, tempVal, compVal)
				}
			}
		} else {
			fmt.Printf("%v is missing from compare file! Value in template is %v\n", key, tempVal)
		}
	}

}

// func askUserInput(path string, templateValue interface{}, compareValue interface{}) interface{} {
// 	retVal := interface{}(nil)

// 	if compareValue == nil {
// 	}

// 	return retVal
// }

// func askUserForArray() []interface{} {
// 	retVal := []interface{}(nil)

// 	fmt.Print("Enter a value to add to the array: ")
// 	un, _ := reader.ReadString('\n')
// 	un = strings.TrimSpace(un)
// 	strLen := utf8.RuneCountInString(un)
// 	for strLen < 3 {
// 		fmt.Printf("\nUsername must be 3 characters! You entered %v.\nEnter your username: ", strconv.Itoa(strLen))
// 		un, _ = reader.ReadString('\n')
// 		un = strings.TrimSpace(un)
// 		strLen = utf8.RuneCountInString(un)
// 	}

// 	return retVal
// }

func flattenJSON(data map[string]*gabs.Container, path string) map[string]interface{} {
	retVals := make(map[string]interface{})
	var recVals map[string]interface{}
	var fullPath string
	for key, child := range data {

		//get the path in string form
		if len(path) < 1 {
			fullPath = key
		} else {
			fullPath = fmt.Sprintf("%v.%v", path, key)
		}

		children, _ := child.ChildrenMap()
		if len(children) > 1 {
			recVals = flattenJSON(children, fullPath)
			for k, v := range recVals {
				retVals[k] = v
			}
		} else {
			switch typeVal := child.Data().(type) {
			case string:
				retVals[fullPath] = typeVal
			case float64:
				retVals[fullPath] = typeVal
			case []interface{}:
				retVals[fullPath] = typeVal
			case map[string]interface{}:
				//this is another json object so send back through
				recVals = flattenJSON(children, fullPath)
				for k, v := range recVals {
					retVals[k] = v
				}
			case bool:
				retVals[fullPath] = typeVal
			default:
				//we didnt find it? Haven't hit this yet
				retVals[fullPath] = typeVal
			}
		}
	}

	return retVals
}
