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
	"wombat/core/util"
)

func Init() {
	go input.run()
}

var input = func() inputManager {
	var i inputManager
	i.quitListeners = make([]func(), 0)
	i.addQuitChan = make(chan func())
	i.removeQuitChan = make(chan func())
	i.mouse = newClicker()
	i.keyboard = newKeyboard()
	return i
}()
//Adds a function to listen for quit requests.
func AddQuit(f func()) {
	input.addQuitChan <- f
}

//Removes a function that listens for quit requests.
func RemoveQuit(f func()) {
	input.removeQuitChan <- f
}

//Adds a function to listen for the clicking of a mouse button.
func AddMouseClick(f func(byte, util.Point)) {
	input.mouse.addClickChan <- f
}

//Removes a function to listen for the pressing of a mouse button.
func RemoveMouseClick(f func(byte, util.Point)) {
	input.mouse.removeClickChan <- f
}

//Adds a function to listen for the pressing of a mouse button.
func AddMousePress(f func(byte, util.Point)) {
	input.mouse.addDownChan <- f
}

//Removes a function to listen for the pressing of a mouse button.
func RemoveMousePress(f func(byte, util.Point)) {
	input.mouse.removeDownChan <- f
}

//Adds a function to listen for the releasing of a mouse button.
func AddMouseRelease(f func(byte, util.Point)) {
	input.mouse.addUpChan <- f
}

//Removes a function to listen for the releasing of a mouse button.
func RemoveMouseRelease(f func(byte, util.Point)) {
	input.mouse.removeUpChan <- f
}

//keyboard listening

//Adds a function to listen for the clicking of a key.
func AddKeyType(f func(KeyEvent)) {
	input.keyboard.addTypeChan <- f
}

//Removes a function to listen for the pressing of a key.
func RemoveKeyType(f func(KeyEvent)) {
	input.keyboard.removeTypeChan <- f
}

//Adds a function to listen for the pressing of a key.
func AddKeyPress(f func(KeyEvent)) {
	input.keyboard.addDownChan <- f
}

//Removes a function to listen for the pressing of a key.
func RemoveKeyPress(f func(KeyEvent)) {
	input.keyboard.removeDownChan <- f
}

//Adds a function to listen for the releasing of a key.
func AddKeyRelease(f func(KeyEvent)) {
	input.keyboard.addUpChan <- f
}

//Removes a function to listen for the releasing of a key.
func RemoveKeyRelease(f func(KeyEvent)) {
	input.keyboard.removeUpChan <- f
}

type inputManager struct {
	quitListeners  []func()
	addQuitChan    chan func()
	removeQuitChan chan func()
	mouse          clicker
	keyboard       keyboard
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
		for loop := true; loop; {
			select {
			case q := <-me.addQuitChan:
				me.quitListeners = append(me.quitListeners, q)
			case r := <-me.removeQuitChan:
				l := me.quitListeners
				var i int
				for i, _ = range l {
					if l[i] == r {
						break
					}
				}

				for i := i; i+1 < len(l); i++ {
					l[i] = l[i+1]
				}
			default:
				loop = false
			}
		}
		time.Sleep(6000000)
	}
}
