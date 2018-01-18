// Copyright (C) 2012-2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package user_agent

import (
 "regexp"
 "strings"
)

// A struct containing all the information that we might be
// interested from the browser.
type Browser struct {
 engine         string
 engine_version string
 name           string
 version        string
}

//D0054
//ulapph
//undefined UserAgent error 1/7/2018
// The UserAgent struct contains all the info that can be extracted
// from the User-Agent string.
type UserAgent struct {
	ua           string
	mozilla      string
	platform     string
	os           string
	localization string
	browser      Browser
	bot          bool
	mobile       bool
	undecided    bool
}

// Internal: extract all the information that we can get from the User-Agent
// string about the browser and update the receiver with this information.
//
// The function receives just one argument "sections", that contains the
// sections from the User-Agent string after being parsed.
func (p *UserAgent) detectBrowser(ua string, sections []section) {
 slen := len(sections)

 if sections[0].name == "Opera" {
  p.mozilla = ""
  p.browser.name = "Opera"
  p.browser.version = sections[0].version
  p.browser.engine = "Presto"
  if slen > 1 {
   p.browser.engine_version = sections[1].version
  }
 } else if slen > 1 {
  engine := sections[1]
  p.browser.engine = engine.name
  p.browser.engine_version = engine.version
  if slen > 2 {
   p.browser.version = sections[2].version
   if engine.name == "AppleWebKit" {
    if sections[slen-1].name == "OPR" {
     p.browser.name = "Opera"
     p.browser.version = sections[slen-1].version
    } else if sections[2].name == "Chrome" {
     p.browser.name = "Chrome"
    } else {
     p.browser.name = "Safari"
    }
   } else if engine.name == "Gecko" {
    p.browser.name = sections[2].name
   } else if engine.name == "like" && sections[2].name == "Gecko" {
    // This is the new user agent from Internet Explorer 11.
    p.browser.engine = "Trident"
    p.browser.name = "Internet Explorer"
    reg, _ := regexp.Compile("^rv:(.+)$")
    for _, c := range sections[0].comment {
     version := reg.FindStringSubmatch(c)
     if len(version) > 0 {
      p.browser.version = version[1]
      return
     }
    }
    p.browser.version = ""
   }
  }
 } else if slen == 1 && len(sections[0].comment) > 1 {
  comment := sections[0].comment
  if comment[0] == "compatible" && strings.HasPrefix(comment[1], "MSIE") {
   p.browser.engine = "Trident"
   p.browser.name = "Internet Explorer"
   // The MSIE version may be reported as the compatibility version.
   // For IE 8 through 10, the Trident token is more accurate.
   // http://msdn.microsoft.com/en-us/library/ie/ms537503(v=vs.85).aspx#VerToken
   for _, v := range comment {
    if strings.HasPrefix(v, "Trident/") {
     switch v[8:] {
     case "4.0":
      p.browser.version = "8.0"
     case "5.0":
      p.browser.version = "9.0"
     case "6.0":
      p.browser.version = "10.0"
	 //ulapph
     case "7.0":
      p.browser.version = "11.0"
     }
     break
    }
   }
   // If the Trident token is not provided, fall back to MSIE token.
   if p.browser.version == "" {
    p.browser.version = strings.TrimSpace(comment[1][4:])
   }
  }
 }
 //ulapph
 if strings.Index(ua, "Edge/") != -1 {
   p.browser.name = "Microsoft Edge"
 }
}

// Public: get the info from the browser's rendering engine.
//
// Returns two strings. The first string is the name of the engine and the
// second one is the version of the engine.
func (p *UserAgent) Engine() (string, string) {
 return p.browser.engine, p.browser.engine_version
}

// Public: get the info from the browser itself.
//
// Returns two strings. The first string is the name of the browser and the
// second one is the version of the browser.
func (p *UserAgent) Browser() (string, string) {
 return p.browser.name, p.browser.version
}
