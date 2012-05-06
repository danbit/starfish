/*
   Copyright 2011-2012 gtalent2@gmail.com

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
package graphics

/*
#cgo LDFLAGS: -lSDL -lSDL_image -lSDL_gfx
#include "SDL/SDL.h"
#include "SDL/SDL_rotozoom.h"
#include "SDL/SDL_image.h"

*/
import "C"

import (
	"encoding/json"
	"github.com/gtalent/starfish/util"
)

type imageLabel struct {
	Str      string
	FilePath bool
}

type imageKey struct {
	Label  imageLabel
	Angle  float64
	Width  int
	Height int
}

func (me *imageKey) String() string {
	str, _ := json.Marshal(me)
	return string(str)
}

var images = newFlyweight(
	func(me *flyweight, path key) interface{} {
		key := path.(*imageKey)
		var i, tmp *C.SDL_Surface
		var cleanup func()
		if key.Label.FilePath {
			tmp = C.IMG_Load(C.CString(key.Label.Str))
			i = C.SDL_DisplayFormatAlpha(tmp)
			C.SDL_FreeSurface(tmp)
			tmp = i
			cleanup = func() { C.SDL_FreeSurface(tmp) }
		} else {
			var k imageKey
			json.Unmarshal([]byte(key.Label.Str), &k)
			i = me.checkout(&k).(*C.SDL_Surface)
			cleanup = func() {}
		}
		var w, h int
		if key.Width == -1 {
			w = int(i.w)
		} else {
			w = key.Width
		}
		if key.Height == -1 {
			h = int(i.h)
		} else {
			h = key.Height
		}
		if (i != nil) && (w != int(i.w) || h != int(i.h) || key.Angle != 0) {
			i = resizeAngleOf(i, key.Angle, w, h)
			cleanup()
		}
		return i
	},
	func(me *flyweight, path key, img interface{}) {
		i := img.(*C.SDL_Surface)
		C.SDL_FreeSurface(i)
	})

type Image struct {
	img *C.SDL_Surface
	key imageKey
}

//Loads the image at the given path, or nil if the image was not found.
func LoadImage(path string) (img *Image) {
	var key imageKey
	key.Label.FilePath = true
	key.Label.Str = path
	key.Width = -1
	key.Height = -1
	i := images.checkout(&key).(*C.SDL_Surface)
	img = new(Image)
	img.img = i
	img.key = key
	return
}

//Loads the image at the given path, or nil if the image was not found.
func LoadImageSize(path string, width, height int) (img *Image) {
	var key imageKey
	key.Label.FilePath = true
	key.Label.Str = path
	key.Width = width
	key.Height = height
	i := images.checkout(&key).(*C.SDL_Surface)
	img = new(Image)
	img.img = i
	img.key = key
	return
}

//Returns the width of the image.
func (me *Image) Width() int {
	return int(me.img.w)
}

//Returns the height of the image.
func (me *Image) Height() int {
	return int(me.img.h)
}

//Returns a util.Size object representing the size of this Image.
func (me *Image) Size() util.Size {
	var s util.Size
	s.Width = me.Width()
	s.Height = me.Height()
	return s
}

//Returns a unique string that can be used to identify the values of this Image.
func (me *Image) String() string {
	return me.key.String()
}

//Returns the path to the image on the disk.
func (me *Image) Path() string {
	return me.key.Label.Str
}

//Nils this image and lets the resource manager know this object is no longer using the image data.
func (me *Image) Free() {
	images.checkin(&me.key)
	me.img = nil
	me.key.Label.Str = ""
}

func resizeAngleOf(img *C.SDL_Surface, angle float64, width, height int) *C.SDL_Surface {
	if img.w == 0 || img.h == 0 {
		return nil
	}
	xstretch := C.double(float64(width) / float64(img.w))
	ystretch := C.double(float64(height) / float64(img.h))
	retval := C.rotozoomSurfaceXY(img, C.double(angle), xstretch, ystretch, 1)
	return retval
}
