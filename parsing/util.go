package parsing

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func marshError(input interface{}, stage string, err error) {
	if err != nil {
		fmt.Println(input)
		fmt.Println(stage)
		log.Fatal("Remashalling error! ", err)
	}
}

// Remarshal: Converts an interface back into a Keyvalue type through re-marshalling
func Remarshal(input interface{}) Keyvalue {
	var back Keyvalue
	out, e := json.Marshal(input)
	marshError(input, "Marshal", e)
	e = json.Unmarshal([]byte(out), &back)
	marshError(input, "Unmarshal", e)
	return back
}

// Slicer: Creates an array of key names given a Keyvalue map
func Slicer(input Keyvalue) []string {
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

// PathFormatter: Given an array, construct it into a jmespath expression (string with . separator)
func PathFormatter(input []string) string {
	var r string
	for i := range input {
		if i == (len(input) - 1) {
			r = r + input[i]
		} else {
			r = r + input[i] + "."
		}
	}
	return r
}

// IndexOf: Finds index of an object in a given array
func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}

// UnorderedKeyMatch: Returns a bool dependant on all 'keys' in a map matching.
func UnorderedKeyMatch(o Keyvalue, m Keyvalue) bool {
	istanbool := true
	oSlice := Slicer(o)
	mSlice := Slicer(m)
	for k := range oSlice {
		val := IndexOf(mSlice, oSlice[k])
		if val == -1 {
			istanbool = false
		}
	}

	for k := range mSlice {
		val := IndexOf(oSlice, mSlice[k])
		if val == -1 {
			istanbool = false
		}
	}
	return istanbool
}

// SliceIndex: Adds an 'index' value to the last string in the slice, used for the 'path' to handle arrays.
func SliceIndex(i int, path []string) []string {

	nPath := make([]string, len(path))
	copy(nPath, path)
	iter := len(nPath) - 1
	nPath[iter] = nPath[iter] + "[" + strconv.Itoa(i) + "]"
	return nPath
}


func MatchAny(compare interface{}, compareSlice []interface{}) bool {
	for i := range compareSlice {
		if reflect.DeepEqual(compare, compareSlice[i]) {
			return true
		}
	}
	return false
}

// DoMapArrayKeysMatch: Uses 'UnorderedKeyMatch' to return a bool for two interfaces if they're both maps
func DoMapArrayKeysMatch(o interface{}, m interface{}) bool {
	if reflect.TypeOf(o).Kind() == reflect.Map && reflect.TypeOf(m).Kind() == reflect.Map {
		return UnorderedKeyMatch(Remarshal(o), Remarshal(m))
	}
	return false
}
