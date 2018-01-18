// Copyright (C) 2012-2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

// Package user_agent implements an HTTP User Agent string parser. It defines
// the type UserAgent that contains all the information from the parsed string.
// It also implements the Parse function and getters for all the relevant
// information that has been extracted from a parsed User Agent string.
package user_agent

import (
 "regexp"
 "strings"
)

// A section contains the name of the product, its version and
// an optional comment.
type section struct {
	name    string
	version string
	comment []string
}

//new
func New(ua string) *UserAgent {
 o := &UserAgent{}
 o.Parse(ua)
 return o
}

// Public: parse the given User-Agent string. After calling this function, the
// receiver will be setted up with all the information that we've extracted.
func (p *UserAgent) Parse(ua string) {
 var sections []section

 p.mobile = false
 for index, limit := 0, len(ua); index < limit; {
  s := parseSection(ua, &index)
  if !p.mobile && s.name == "Mobile" {
   p.mobile = true
  }
  sections = append(sections, s)
 }

 if len(sections) > 0 {
  p.mozilla = sections[0].version
  if !p.bot {
   for _, v := range sections {
    p.checkBot(v.comment)
   }
   if !p.bot {
    p.detectBrowser(ua, sections)
    p.detectOS(sections[0])
   }
  }
 }
}

// Internal: check if we're dealing with a Bot.
func (p *UserAgent) checkBot(comment []string) {
 reg, _ := regexp.Compile("(?i)bot")
 for _, v := range comment {
  if reg.Match([]byte(v)) {
   p.bot = true
   return
  }
 }
}

//D0054
//1/7/2018
// Parse a section. A section is typically formatted as follows
// "Name/Version (comment)". Both, the comment and the version are optional.
//
// The first argument is the user agent string being parsed. The second
// argument is a reference pointing to the current index of the user agent
// string.
//
// Returns a section containing the information that we could extract
// from the last parsed section.
func parseSection(ua string, index *int) (s section) {
	buffer := readUntil(ua, index, ' ', false)

	s.name, s.version = parseProduct(buffer)
	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer = readUntil(ua, index, ')', true)
		s.comment = strings.Split(string(buffer), "; ")
		*index++
	}
	return s
}

//D0054
//1/7/2018
// Returns an array of bytes containing what has been read.
func readUntil(ua string, index *int, delimiter byte, cat bool) []byte {
	var buffer []byte

	i := *index
	catalan := 0
	for ; i < len(ua); i = i + 1 {
		if ua[i] == delimiter {
			if catalan == 0 {
				*index = i + 1
				return buffer
			}
			catalan--
		} else if cat && ua[i] == '(' {
			catalan++
		}
		buffer = append(buffer, ua[i])
	}
	*index = i + 1
	return buffer
}

//D0054
//1/7/2018
// It returns two strings. The first string is the name of the product and the
// second string contains the version of the product.
func parseProduct(product []byte) (string, string) {
	prod := strings.SplitN(string(product), "/", 2)
	if len(prod) == 2 {
		return prod[0], prod[1]
	}
	return string(product), ""
}

// Public: get the mozilla version (it's how the User Agent string begins:
// "Mozilla/5.0 ...", unless we're dealing with Opera, of course).
// Returns a string containing the mozilla version.
func (p *UserAgent) Mozilla() string {
 return p.mozilla
}

// Returns true if it's a bot, false otherwise.
func (p *UserAgent) Bot() bool {
 return p.bot
}

// Returns true if it's a mobile device, false otherwise.
func (p *UserAgent) Mobile() bool {
 return p.mobile
}
