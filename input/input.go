/*
   Copyright 2011 gtalent2@gmail.com

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
package input

import (
	"time"
	"sdl"
)

func Init() {
	input.run()
}

type inputManager struct {
	quitListeners []func()
	mouse         clicker
}

var input = func() inputManager {
	var i inputManager
	i.quitListeners = make([]func(), 0)
	i.mouse = newClicker()
	return i
}()

//Adds a function to listn for the pressing of a mouse button.
func (me *inputManager) AddMouseDown(f func(byte)) {
	me.mouse.addDownChan <- f
}

//Removes a function to listn for the pressing of a mouse button.
func (me *inputManager) RemoveMouseDown(f func(byte)) {
	me.mouse.removeDownChan <- f
}

//Adds a function to listn for the releasing of a mouse button.
func (me *inputManager) AddMouseUp(f func(byte)) {
	me.mouse.addUpChan <- f
}

//Removes a function to listn for the releasing of a mouse button.
func (me *inputManager) RemoveMouseUp(f func(byte)) {
	me.mouse.removeUpChan <- f
}

func (me *inputManager) run() {
	for {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch et := e.(type) {
			case *sdl.QuitEvent:
				for _, a := range me.quitListeners {
					go a()
				}
			case *sdl.MouseButtonEvent:
				me.mouse.input <- et
			}
		}
		time.Sleep(6000000)
	}
}
