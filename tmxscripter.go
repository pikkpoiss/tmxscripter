// Copyright 2015 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/kurrik/fauxfile"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
)

type TmxScripter struct {
	fs         fauxfile.Filesystem
	InputPath  string
	OutputPath string
	ScriptPath string
}

func NewTmxScripter(fs fauxfile.Filesystem) *TmxScripter {
	return &TmxScripter{
		fs: fs,
	}
}

func (s *TmxScripter) loadMap() (m *tmxgo.Map, err error) {
	var (
		f         fauxfile.File
		inputFile []byte
	)
	if f, err = s.fs.Open(s.InputPath); err != nil {
		err = fmt.Errorf("Could not open input file: %v", err)
		return
	}
	defer f.Close()
	if inputFile, err = ioutil.ReadAll(f); err != nil {
		err = fmt.Errorf("Could not read input file: %v", err)
		return
	}
	if m, err = tmxgo.ParseMapString(string(inputFile)); err != nil {
		err = fmt.Errorf("Could not parse map file: %v", err)
		return
	}
	return
}

func (s *TmxScripter) saveMap(m *tmxgo.Map) (err error) {
	var (
		f          fauxfile.File
		serialized string
	)
	if serialized, err = m.Serialize(); err != nil {
		err = fmt.Errorf("Could not reserialize map: %v", err)
		return
	}
	if f, err = s.fs.Create(s.OutputPath); err != nil {
		err = fmt.Errorf("Could not open output file: %v", err)
		return
	}
	defer f.Close()
	if _, err = f.Write([]byte(serialized)); err != nil {
		err = fmt.Errorf("Could not write output file: %v", err)
		return
	}
	return
}

func (s *TmxScripter) Run() (err error) {
	var (
		m *tmxgo.Map
	)
	if m, err = s.loadMap(); err != nil {
		return
	}
	if err = s.saveMap(m); err != nil {
		return
	}
	return fmt.Errorf("Not implemented")
}
