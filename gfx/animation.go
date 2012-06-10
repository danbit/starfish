/*
   Copyright 2011-2012 starfish authors

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
package gfx

import (
	"strconv"
	"time"
)

var animTicker = time.Now().UnixNano()
var tickerInterval int64 = 100000000

//Sets how often the clock that Animations update based on is updated.
func SetAnimTickInterval(ms int64) {
	tickerInterval = ms
}

func startAnimTick() {
	go func() {
		for {
			select {
			case <-kill:
				return
			default:
				animTicker = time.Now().UnixNano()
				time.Sleep(time.Duration(tickerInterval))
			}
		}
	}()
}


//A type to automatically flip through a series of images.
type Animation struct {
	interval   int64
	lastUpdate int64
	slide      int
	images     []*Image
}

func NewAnimation(interval int) *Animation {
	a := new(Animation)
	a.SetInterval(interval)
	return a
}

//Returns a string that can be used to identify the values of this Animation.
func (me *Animation) String() string {
	retval := strconv.FormatInt(me.interval, 10)
	for _, i := range me.images {
		retval += "\n" + i.String()
	}
	return retval
}

//Sets the number of milliseconds per image.
func (me *Animation) SetInterval(ms int) {
	me.interval = int64(ms) * 1000000
}

//Gets the current image.
func (me *Animation) GetImage() *Image {
	if me.images == nil {
		return nil
	}
	if t := animTicker; t-me.lastUpdate >= me.interval {
		me.slide += int((t - me.lastUpdate) / me.interval)
		me.slide %= len(me.images)
		me.lastUpdate = t
	}
	return me.images[me.slide]
}

//Returns the image at the given index.
func (me *Animation) At(i int) *Image {
	return me.images[i]
}

//Returns the number of images in this Animation.
func (me *Animation) Size() int {
	return len(me.images)
}

func (me *Animation) LoadImage(path string) {
	if i := LoadImage(path); i != nil {
		me.images = append(me.images, i)
	}
}

func (me *Animation) LoadImageSize(path string, width, height int) {
	if i := LoadImageSize(path, width, height); i != nil {
		me.images = append(me.images, i)
	}
}

//Frees this Animations images, rendering it useless.
func (me *Animation) Free() {
	for _, a := range me.images {
		a.Free()
	}
}
