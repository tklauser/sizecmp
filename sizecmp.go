// Copyright 2020 Tobias Klauser. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on https://github.com/rsc/tmp/blob/master/sizecmp/main.go

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: sizecmp old.txt new.txt")
		os.Exit(1)
	}

	sz1 := readSize(os.Args[1])
	sz2 := readSize(os.Args[2])

	var keys []string
	for k := range sz1 {
		if _, ok := sz2[k]; ok {
			keys = append(keys, k)
		} else {
			fmt.Printf("binary %q found in %q but not in %q", k, os.Args[2], os.Args[1])
		}
	}

	for k := range sz2 {
		if _, ok := sz1[k]; !ok {
			fmt.Printf("binary %q found in %q but not in %q", k, os.Args[1], os.Args[2])
		}
	}
	sort.Strings(keys)

	for _, kk := range keys {
		var skeys []string
		for k, _ := range sz1[kk] {
			skeys = append(skeys, k)
		}
		sort.Strings(skeys)

		fmt.Printf("== %s ==\n", kk)

		var total1, total2 int64
		for _, k := range skeys {
			s1 := sz1[kk][k]
			s2 := sz2[kk][k]
			fmt.Printf("%-30s %11d %11d %+11d\n", k, s1, s2, s2-s1)
			total1 += s1
			total2 += s2
		}
		fmt.Printf("%30s %11d %11d %+11d\n", "total", total1, total2, total2-total1)
	}
}

// readSize reads a file containing the output of the `size` command.
//
// Example output:
//    text	   data	    bss	    dec	    hex	filename
// 55876285	 893753	7752192	64522230	3d887f6	daemon/cilium-agent
// 57015715	 610160	 211608	57837483	37287ab	operator/cilium-operator
// 35449857	 486256	 192536	36128649	2274789	plugins/cilium-cni/cilium-cni
func readSize(file string) map[string]map[string]int64 {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var header []string
	objs := make(map[string]map[string]int64)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 6 {
			log.Fatal("invalid file format in %q", file)
		}

		if fields[0] == "text" {
			header = fields
			continue
		}

		sizes := make(map[string]int64)
		for i, field := range fields {
			if i < 5 {
				v, err := strconv.ParseInt(field, 16, 64)
				if err != nil {
					log.Fatal(err)
				}
				sizes[header[i]] = v
			} else {
				objs[field] = sizes
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return objs
}
