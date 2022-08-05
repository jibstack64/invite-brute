package brute

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// A default code generator with default values.
var DefCodeGenerator = &CodeGenerator{
	Chars: DefaultChars[:],
}

var (
	// The "default" array of Discord invite code characters.
	DefaultChars = [62]rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z',
	}
)

// Contains functions and values for generating Discord invite codes.
type CodeGenerator struct {
	// The slice of characters to be used in generating codes.
	Chars []rune
}

/*
	Generates possible Discord invite codes.

`n` is the number of codes to be generated.

Coupled, `minLength` and `maxLength` represent the lowest and highest
number of characters within a code. The resulting length is randomly
generated, but is ensured to be within the range of these integers.
*/
func (c *CodeGenerator) GenerateCodes(n int, minLength int, maxLength int) (fc *[]string) {
	codes := make([]string, 0, n)
	for i := 0; i < n; i++ {
		code := ""
		for ii := 0; ii < rand.Intn(maxLength-minLength)+minLength; ii++ {
			// Get a random rune, convert it to string and append it to the code string
			code += strings.ReplaceAll(strconv.QuoteRune(c.Chars[rand.Intn(len(c.Chars))]), "'", "")
		}
		// Finally, append the code to the code slice
		codes = append(codes, code)
	}
	fc = &codes
	return
}

/*
	Writes all codes provided to the file specified by the given path.

This function only suits situations where the program needs to ensure
that the provided codes are not lost.
*/
func (c *CodeGenerator) WriteToFile(path string, codes ...string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return
	} else {
		defer f.Close()
	}
	// Create string of all codes for writing
	cs := ""
	for _, code := range codes {
		cs += code + "\n"
	}
	// Encode the string and write out
	_, err = f.Write([]byte(cs))
	return
}

/*
Reads all codes from a file.
*/
func (c *CodeGenerator) ReadFromFile(path string) (fc *[]string, err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0664)
	if err != nil {
		return
	} else {
		defer f.Close()
		// Read all file bytes line by line
		fileScanner := bufio.NewScanner(f)
		fileScanner.Split(bufio.ScanLines)
		// Temporary fc
		c := make([]string, 0)
		for fileScanner.Scan() {
			c = append(c, fileScanner.Text())
		}
		// Create pointer for string slice
		fc = &c
		return
	}
}
