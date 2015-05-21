/*
	(c) Copyright 2015 Vlad Didenko

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package persist

import (
	"errors"
	"io"
)

type StringReadWriteSeeker struct {
	s   string
	pos int
}

const maxint = int64(^uint(0) >> 1)

func (sw *StringReadWriteSeeker) Write(p []byte) (n int, err error) {
	sw.s = sw.s[:sw.pos] + string(p)
	return len(p), nil
}

func (sw *StringReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case 0:
		newOffset = offset
	case 1:
		newOffset = int64(sw.pos) + offset
	case 2:
		newOffset = int64(len(sw.s)) + offset
	}

	switch {
	case newOffset < 0:
		return 0, errors.New("Negative offset is not allowed")
	case newOffset > int64(len(sw.s)):
		return newOffset, errors.New("Requested offset is beyond the string length")
	default:
		sw.pos = int(newOffset)
		return newOffset, nil
	}
}

func (sw *StringReadWriteSeeker) Read(p []byte) (n int, err error) {
	is, ip := 0, 0
	next := false

	lp := len(p) - 1
	ls := len(sw.s) - 1

	if sw.pos > ls {
		return 0, io.EOF
	}

	for is, ip, next = sw.pos, 0, (is < ls && ip < lp); next; is, ip, next = is+1, ip+1, (is < ls && ip < lp) {
		p[ip] = sw.s[is]
	}
	sw.pos += ip
	return ip, nil
}

func (sw *StringReadWriteSeeker) Get() string {
	return sw.s
}

func (sw *StringReadWriteSeeker) Set(s string) {
	sw.pos = 0
	sw.s = s
}
