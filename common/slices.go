package common

import (
	"strings"
)

type StringSlice struct {
	values []string
}

func NewStringSlice(values ...string) *StringSlice {
	return &StringSlice{
		values: values,
	}
}

func (s *StringSlice) Join(separator string) *String {
	return &String{
		value: strings.Join(s.values, separator),
	}
}

func (s *StringSlice) Values() []string {
	return s.values
}

func (s *StringSlice) Contains(value string) bool {
	for _, element := range s.Values() {
		if element == value {
			return true
		}
	}
	return false
}

func (s *StringSlice) AnyPrefix(prefix string) bool {
	for _, element := range s.Values() {
		if strings.HasPrefix(element, prefix) {
			return true
		}
	}
	return false
}
