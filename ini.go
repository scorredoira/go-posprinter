// Package ini allows to read ini files
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Path     string
	Sections []string
	Values   map[string]string
}

func NewConfig(path string) (*Config, error) {
	c, err := Open(path)
	if err != nil {
		return nil, err
	}
	if c != nil {
		c.Path = path
		return c, nil
	}
	c = &Config{
		Path:   path,
		Values: make(map[string]string),
	}

	return c, nil
}

// Get returns a configuration value.
func (c *Config) Get(keys ...string) string {
	key := strings.Join(keys, ".")
	return c.Values[key]
}

// Set assigns a configuration value
func (c *Config) Set(key, value string) {
	c.Values[key] = value
}

// String returns the value if present or defaultValue
func (c *Config) String(key string, defaultValue string) string {
	v, ok := c.Values[key]
	if ok {
		return v
	}

	return defaultValue
}

// Bool returns the value if present or defaultValue
func (c *Config) Bool(key string, defaultValue bool) bool {
	v, ok := c.Values[key]
	if ok {
		return v == "true"
	}

	return defaultValue
}

// Int returns the value if present or defaultValue
func (c *Config) Int(key string, defaultValue int) int {
	v, ok := c.Values[key]
	if ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// File parses a ini file
func Open(f string) (*Config, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return Parse(string(b))
}

func (c *Config) Save() error {
	var buf []byte

	for k, v := range c.Values {
		if strings.ContainsRune(k, '.') {
			return fmt.Errorf("Section save is not implemented yet")
		}
		buf = append(buf, []byte(k+"="+v+"\n")...)
	}

	return ioutil.WriteFile(c.Path, buf, 0644)
}

// File parses a ini string
func Parse(s string) (*Config, error) {
	var sections []string
	values := make(map[string]string)
	currentSection := ""

	for l, line := range strings.Split(s, "\n") {
		line = strings.Trim(line, " \t\r")

		// ignore empty lines or comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// parse section
		if line[0] == '[' {
			if line[len(line)-1] != ']' {
				return nil, fmt.Errorf("Unbound section in line %d: %s", l, line)
			}

			currentSection = line[1 : len(line)-1]
			sections = append(sections, currentSection)
			continue
		}

		// parse value
		i := strings.Index(line, "=")
		if i == -1 {
			return nil, fmt.Errorf("Invalid value in line %d: %s", l, line)
		}

		key := strings.Trim(line[:i], " \t")
		if currentSection != "" {
			key = currentSection + "." + key
		}

		value := strings.Trim(line[i+1:], " \t")
		values[key] = value
	}

	return &Config{Sections: sections, Values: values}, nil
}
