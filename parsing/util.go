package parsing

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"runtime/debug"
	"golang.org/x/text/unicode/rangetable"
	"regexp"
)

func marshError(input interface{}, stage string, err error) {
	if err != nil {
		fmt.Println(input)
		fmt.Println(stage)
		debug.PrintStack()
		log.Fatal("Remashalling error! ", err)

	}
}

// Remarshal deprecated
func Remarshal(input interface{}) KeyValue {
	// This is just a nasty type conversions, marshals an interface and then back into our Keyvalue map type
	var back KeyValue
	out, e := json.Marshal(input)
	marshError(input, "Marshal", e)
	e = json.Unmarshal([]byte(out), &back)
	marshError(input, "Unmarshal", e)
	return back
}


// GetSliceOfKeys creates a slice of keys from an object
func GetSliceOfKeys(input KeyValue) []string {
	// Creates an array of key names given a Keyvalue map
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

// CreatePath Given an array, construct it into a jmespath expression (string with . separator)
func CreatePath(input []string) string {
	var r string


	// characters to escape
	escapeChars := ".-:"

	// iterate over path slice to construct string path
	for i,str := range input {
		var indexStr string
		wrappedReg := regexp.MustCompile("^\".*\"(\\[[\\d]\\])*$")
		indexReg := regexp.MustCompile("\\[[\\d]\\]+")


		// if string has index values iterate them all into string
		index := indexReg.FindAllString(str, -1)

		for i := range index {
			indexStr = indexStr + index[i]
		}

		raw := indexReg.ReplaceAllString(str, "")

		// Escape a . in string name for parsing later
		if !wrappedReg.MatchString(str) && strings.ContainsAny(str, escapeChars) {
			str = strconv.Quote(raw)
			str = str + indexStr
			//str = "\"" + raw + "\"" + indexStr
		}

		if i == (len(input) - 1) {
			r = r + str
		} else {
			r = r + str + "."
		}

	}
	return r
}

// IndexOf Finds index of an object in a given array
func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}

// SliceIndexOf find index of item from in slice
func SliceIndexOf(item interface{}, list []interface{}) int {
	for i := 0; i < len(list); i++ {
		if list[i] != nil {
			if reflect.DeepEqual(item, list[i]) {
				return i
			}
		}
	}
	return -1
}


// UnorderedKeyMatch Returns a bool dependant on all 'keys' in a map matching.
func UnorderedKeyMatch(o map[string]interface{}, m map[string]interface{}) bool {
	istanbool := true
	oSlice := GetSliceOfKeys(o)
	mSlice := GetSliceOfKeys(m)
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

// SliceIndex Adds an 'index' value to the last string in the slice, used for the 'path' to handle arrays.
func SliceIndex(i int, path []string) []string {

	nPath := make([]string, len(path))
	copy(nPath, path)
	iter := len(nPath) - 1
	nPath[iter] = nPath[iter] + "[" + strconv.Itoa(i) + "]"
	return nPath
}

// MatchAny check if an object exists anywhere in a slice
func MatchAny(compare interface{}, compareSlice []interface{}) bool {
	for i := range compareSlice {
		if reflect.DeepEqual(compare, compareSlice[i]) {
			return true
		}
	}
	return false
}


// MapMatchAny check if map exists in larger map
func MapMatchAny(a map[string]interface{}, b map[string]interface{}) bool {
	for k,v  :=  range b {
		c := map[string]interface{}{
			k:v,
		}
		if reflect.DeepEqual(a, c) {
			return true
		}
	}
	return false
}



// DoMapArrayKeysMatch Uses 'UnorderedKeyMatch' to return a bool for two interfaces if they're both maps
func DoMapArrayKeysMatch(o interface{}, m interface{}) bool {
	if reflect.TypeOf(o).Kind() == reflect.Map && reflect.TypeOf(m).Kind() == reflect.Map {
		return UnorderedKeyMatch(Remarshal(o), Remarshal(m))
	}
	return false
}

// PathSplit Splits up jmespath format path into a slice, will ignore escaped '.' ; opposite of CreatePath
func PathSplit(input string) []string {
	return escape(input)
}

func escape(input string) []string {
	dotRange := rangetable.New(rune('.'))
	old := rune(0)
	f := func(c rune) bool {
		switch {
		case c == old:
			old = rune(0)
			return false
		case old != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			old = c
			return false
		default:
			return unicode.In(c, dotRange)
		}
	}
	return strings.FieldsFunc(input, f)

}

// \ = U+005C
// . = U+002E