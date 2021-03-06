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

package tmxscripter

import (
	"fmt"
	"github.com/kurrik/fauxfile"
	"github.com/kurrik/tmxgo"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"io/ioutil"
	"path"
)

type TmxScripter struct {
	fs         fauxfile.Filesystem
	vm         *otto.Otto
	listeners  map[string][]otto.Value
	InputPath  string
	OutputPath string
	ScriptPath string
}

func NewTmxScripter(fs fauxfile.Filesystem) *TmxScripter {
	return &TmxScripter{
		fs:        fs,
		vm:        otto.New(),
		listeners: map[string][]otto.Value{},
	}
}

func (s *TmxScripter) setAPI() {
	s.vm.Set("addEventListener", func(call otto.FunctionCall) otto.Value {
		var usage = fmt.Errorf("Usage: addEventListener(string, func)")
		if len(call.ArgumentList) != 2 {
			panic(usage)
		}
		if !call.Argument(0).IsString() {
			panic(usage)
		}
		if !call.Argument(1).IsFunction() {
			panic(usage)
		}
		var (
			callbacks []otto.Value
			present   bool
			eventName = call.Argument(0).String()
		)
		if callbacks, present = s.listeners[eventName]; !present {
			s.listeners[eventName] = callbacks
		}
		s.listeners[eventName] = append(s.listeners[eventName], call.Argument(1))
		return otto.NullValue()
	})
	s.vm.Set("readFile", func(call otto.FunctionCall) otto.Value {
		var usage = fmt.Errorf("Usage: readFile(path)")
		if len(call.ArgumentList) != 1 {
			panic(usage)
		}
		if !call.Argument(0).IsString() {
			panic(usage)
		}
		var (
			f         fauxfile.File
			fileBytes []byte
			err       error
			filePath  = call.Argument(0).String()
			fullPath  = path.Join(path.Dir(s.ScriptPath), filePath)
			response  otto.Value
		)
		if f, err = s.fs.Open(fullPath); err != nil {
			panic(err)
		}
		defer f.Close()
		if fileBytes, err = ioutil.ReadAll(f); err != nil {
			panic(err)
		}
		if response, err = otto.ToValue(string(fileBytes)); err != nil {
			panic(err)
		}
		return response
	})
	s.vm.Set("writeFile", func(call otto.FunctionCall) otto.Value {
		var usage = fmt.Errorf("Usage: writeFile(path, data)")
		if len(call.ArgumentList) != 2 {
			panic(usage)
		}
		if !call.Argument(0).IsString() {
			panic(usage)
		}
		if !call.Argument(1).IsString() {
			panic(usage)
		}
		var (
			f        fauxfile.File
			err      error
			filePath = call.Argument(0).String()
			fullPath = path.Join(path.Dir(s.ScriptPath), filePath)
		)
		if f, err = s.fs.Create(fullPath); err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err = f.Write([]byte(call.Argument(1).String())); err != nil {
			panic(err)
		}
		return otto.NullValue()
	})
}

func (s *TmxScripter) loadScript() (err error) {
	var (
		f          fauxfile.File
		scriptFile []byte
		script     *otto.Script
	)
	if f, err = s.fs.Open(s.ScriptPath); err != nil {
		err = fmt.Errorf("Could not open script file: %v", err)
		return
	}
	defer f.Close()
	if scriptFile, err = ioutil.ReadAll(f); err != nil {
		err = fmt.Errorf("Could not read script file: %v", err)
		return
	}
	if script, err = s.vm.Compile("", string(scriptFile)); err != nil {
		err = fmt.Errorf("Could not compile script: %v", err)
		return
	}
	s.setAPI()
	if _, err = s.vm.Run(script); err != nil {
		err = fmt.Errorf("Could not execute script: %v", err)
		return
	}
	return
}

func (s *TmxScripter) convertArgs(args []interface{}) (converted []interface{}, err error) {
	var (
		i   int
		arg interface{}
	)
	converted = make([]interface{}, len(args))
	for i, arg = range args {
		if converted[i], err = s.vm.ToValue(arg); err != nil {
			return
		}
	}
	return
}

func (s *TmxScripter) triggerEvent(eventName string, rawArgs ...interface{}) (err error) {
	var (
		present   bool
		callbacks []otto.Value
		callback  otto.Value
		arguments []interface{}
	)
	if callbacks, present = s.listeners[eventName]; !present {
		return // Not an error to have no listeners.
	}
	if arguments, err = s.convertArgs(rawArgs); err != nil {
		return
	}
	for _, callback = range callbacks {
		if _, err = callback.Call(callback, arguments...); err != nil {
			return
		}
	}
	return
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

func (s *TmxScripter) Validate() (err error) {
	if _, err = s.fs.Stat(s.InputPath); err != nil {
		return fmt.Errorf("Input file %v does not exist", s.InputPath)
	}
	if _, err = s.fs.Stat(s.ScriptPath); err != nil {
		return fmt.Errorf("Script file %v does not exist", s.ScriptPath)
	}
	return
}

func (s *TmxScripter) Run() (err error) {
	var (
		m *tmxgo.Map
	)
	if err = s.Validate(); err != nil {
		return
	}
	if m, err = s.loadMap(); err != nil {
		return
	}
	if err = s.loadScript(); err != nil {
		return
	}
	if err = s.triggerEvent("map", NewScriptableMap(m)); err != nil {
		return
	}
	if err = s.saveMap(m); err != nil {
		return
	}
	return
}
