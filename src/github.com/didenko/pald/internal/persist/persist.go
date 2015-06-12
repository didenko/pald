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
	"io"
	"time"
)

type Dumper interface {
	Dump(w io.Writer) (bytes int, err error)
}

type RWST interface {
	io.ReadWriteSeeker
	Truncate(size int64) error
}

type Loader interface {
	Load(r io.Reader) (err error)
}

func Persist(src Dumper, dst RWST, throttle time.Duration) chan struct{} {

	defer src.Dump(dst)

	knob := make(chan struct{}, 10)

	go func() {

		for {

			last := time.Now()

			select {

			case _, ok := <-knob:
				if !ok {
					break
				}

				if time.Since(last) > throttle {
					dst.Seek(0, 0)
					dst.Truncate(0)
					src.Dump(dst)
				}
			}
		}
	}()

	return knob
}

func Load(src RWST, dst Loader) error {
	src.Seek(0, 0)
	return dst.Load(src)
}
