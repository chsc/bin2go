// Copyright (c) 2011, Christoph Schunk
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the author nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	pkgName   = flag.String("p", "main", "Package name")
	lineLen   = flag.Int("l", 8, "Line length")
	comment   = flag.Bool("c", false, "Line comments")
	not_sized = flag.Bool("z", false, "Use non-sized []byte")
	single    = flag.String("s", "", "Single output filename")
)

func bin2go(ifile, pkgName, bufName string, ofi *os.File, line int, comment, not_sized, use_single bool) error {
	ifi, err := os.Open(ifile)
	if err != nil {
		return err
	}
	defer ifi.Close()

	buffer := make([]byte, line)

	size := "..."
	if not_sized {
		size = ""
	}

	if !use_single {
		fmt.Fprintf(ofi, "// Automatically generated with bin2go: http://github.com/chsc/bin2go\n")
		fmt.Fprintf(ofi, "package %s\n", pkgName)
	}

	fmt.Fprintf(ofi, "\nvar %s = [%s]byte{\n", bufName, size)
	for {
		nRead, err := ifi.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Fprintf(ofi, "\t%0#2x,", buffer[0])
		for _, c := range buffer[1:nRead] {
			fmt.Fprintf(ofi, " %0#2x,", c)
		}
		if comment {
			fmt.Fprintf(ofi, "\t// %s", clean(string(buffer)))
		}
		fmt.Fprint(ofi, "\n")
	}
	fmt.Fprint(ofi, "}\n")

	return nil
}

func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)
}

func main() {
	flag.Parse()

	var ofi *os.File
	var err error
	var use_single = *single != ""
	if use_single {
		if ofi, err = os.Create(*single); err != nil {
			return
		}

		fmt.Fprintf(ofi, "// Automatically generated with bin2go: http://github.com/chsc/bin2go\n")
		fmt.Fprintf(ofi, "package %s\n", *pkgName)

		defer ofi.Close()
	}

	for _, fileName := range flag.Args() {
		if !use_single {
			if ofi, err = os.Create(fileName + ".go"); err != nil {
				return
			}
		}

		bin2go(fileName, *pkgName, clean(fileName), ofi, *lineLen, *comment, *not_sized, use_single)

		if !use_single {
			ofi.Close()
		}
	}
}
