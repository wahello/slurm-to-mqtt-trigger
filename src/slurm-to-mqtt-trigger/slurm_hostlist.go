package main

/*
 * This hostlist expansion has been migrated from the original Python implementation:
 *
 * Hostlist library
 *
 * Copyright (C) 2008-2018
 *                    Kent Engström <kent@nsc.liu.se>,
 *                    Thomas Bellman <bellman@nsc.liu.se>,
 *                    Pär Lindfors <paran@nsc.liu.se> and
 *                    Torbjörn Lönnemark <ketl@nsc.liu.se>,
 *                    National Supercomputer Centre
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301, USA.
 */

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func expandPart(p string) []string {
	var result []string
	var usExpanded []string

	re := regexp.MustCompile(`([^,\[]*)(\[[^\]]*\])?(.*)`)

	if p == "" {
		return result
	}

	sbm := re.FindStringSubmatch(p)
	prefix := sbm[1]
	rangelist := sbm[2]
	rest := sbm[3]

	restExpanded := expandPart(rest)

	if rangelist == "" {
		usExpanded = append(usExpanded, prefix)
	} else {
		usExpanded = expandRangeList(prefix, rangelist[1:len(rangelist)-1])
	}

	result = append(result, usExpanded...)
	result = append(result, restExpanded...)

	return result
}

func expandRangeList(prefix string, rangelist string) []string {
	var result []string

	for _, r := range strings.Split(rangelist, ",") {
		result = append(result, expandRange(prefix, r)...)
	}

	return result
}

func expandRange(prefix string, r string) []string {
	var result []string

	reSingleNumber := regexp.MustCompile(`^[0-9]+$`)
	if reSingleNumber.MatchString(r) {
		result = append(result, prefix+r)
		return result
	}

	reLowHigh := regexp.MustCompile(`^([0-9]+)-([0-9]+)$`)

	if !reLowHigh.MatchString(r) {
		panic("bad range")
	}

	lhm := reLowHigh.FindStringSubmatch(r)

	_low := lhm[1]
	_high := lhm[2]

	low, err := strconv.ParseInt(_low, 10, 64)
	if err != nil {
		panic(err)
	}

	high, err := strconv.ParseInt(_high, 10, 64)
	if err != nil {
		panic(err)
	}

	width := len(_low)
	if high < low {
		fmt.Println(low, high)
		panic("start > stop")
	}

	for i := low; i <= high; i++ {
		result = append(result, fmt.Sprintf("%s%0*d", prefix, width, i))
	}

	return result
}

func expandSLURMHostlist(hostlist string) []string {
	var result []string
	var _result []string
	var bracketLevel int64
	var part string

	for _, c := range hostlist + "," {
		if c == ',' && bracketLevel == 0 {
			if part != "" {
				_result = append(result, expandPart(part)...)
			}
			part = ""
		} else {
			part += string(c)
		}

		if c == '[' {
			bracketLevel++
		} else if c == ']' {
			bracketLevel--
		}

		if bracketLevel > 1 {
			panic("nested brackets")
		} else if bracketLevel < 0 {
			panic("unbalanced brackets")
		}
	}

	if bracketLevel > 0 {
		panic("unbalanced brackets")
	}

	noDupes := make(map[string]interface{})
	for _, h := range _result {
		noDupes[h] = nil
	}

	for key := range noDupes {
		result = append(result, key)
	}

	return result
}
