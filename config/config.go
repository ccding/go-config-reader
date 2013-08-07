// Copyright 2013, Cong Ding. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: Cong Ding <dinggnu@gmail.com>
//
package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var commentPrefix = []string{"//", "#", ";"}

type Config struct {
	filename string
	config   map[string]string
}

func NewConfig(filename string) *Config {
	c := new(Config)
	c.filename = filename
	c.config = make(map[string]string)
	return c
}

func (c *Config) Read() error {
	in, err := os.Open(c.filename)
	if err != nil {
		return err
	}
	defer in.Close()
	scanner := bufio.NewScanner(in)
	line := ""
	section := ""
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		if line == "" {
			sec := checkSection(scanner.Text())
			if sec != "" {
				section = sec + "."
				continue
			}
		}
		if checkComment(scanner.Text()) {
			continue
		}
		line += scanner.Text()
		if strings.HasSuffix(line, "\\") {
			line = line[:len(line)-1]
			continue
		}
		key, value, err := checkLine(line)
		if err != nil {
			return errors.New("WRONG: " + line)
		}
		c.config[section+key] = value
		line = ""
	}
	return nil
}

func (c *Config) Get(key string) string {
	value, ok := c.config[key]
	if !ok {
		return ""
	}
	return value
}

func (c *Config) Set(key string, value string) {
	c.config[key] = value
}

func (c *Config) Add(key string, value string) {
	c.config[key] = value
}

func (c *Config) Del(key string) {
	delete(c.config, key)
}

func (c *Config) Write() {
}

func checkSection(line string) string {
	line = strings.TrimSpace(line)
	lineLen := len(line)
	if lineLen < 2 {
		return ""
	}
	if line[0] == '[' && line[lineLen-1] == ']' {
		return line[1 : lineLen-1]
	}
	return ""
}

func checkLine(line string) (string, string, error) {
	key := ""
	value := ""
	sp := strings.SplitN(line, "=", 2)
	if len(sp) != 2 {
		return key, value, errors.New("WRONG: " + line)
	}
	key = strings.TrimSpace(sp[0])
	value = strings.TrimSpace(sp[1])
	return key, value, nil
}

func checkComment(line string) bool {
	line = strings.TrimSpace(line)
	for p := range commentPrefix {
		if strings.HasPrefix(line, commentPrefix[p]) {
			return true
		}
	}
	return false
}
