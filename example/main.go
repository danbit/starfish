/*
   Copyright 2011-2014 starfish authors

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
package main

import (
	"flag"
	"fmt"
	l "log"
	"os"
	"runtime"
	"runtime/pprof"
	"github.com/gtalent/starfish/gfx"
	"github.com/gtalent/starfish/input"
)


var log = l.New(os.Stdout, "  LOG: starfish example: ", l.Ldate|l.Ltime)
var errlog = l.New(os.Stderr, "ERROR: starfish example: ", l.Ldate|l.Ltime)

func LogOn(on bool) {
	if on {
		log = l.New(os.Stdout, "  LOG: starfish example: ", l.Ldate|l.Ltime)
	} else {
		log = l.New(nil, "  LOG: starfish example: ", l.Ldate|l.Ltime)
	}
}

func ErrLogOn(on bool) {
	if on {
		errlog = l.New(os.Stdout, "ERROR: starfish backend: ", l.Ldate|l.Ltime)
	} else {
		errlog = l.New(nil, "ERROR: starfish backend: ", l.Ldate|l.Ltime)
	}
}

type Drawer struct {
	box  *gfx.Image
	text *gfx.Text
	anim *gfx.Animation
}

func (me *Drawer) init() bool {
	me.anim = gfx.NewAnimation(1000)
	me.box = gfx.LoadImageSize("box.png", 100, 100)
	font := gfx.LoadFont("LiberationSans-Bold.ttf", 32)
	if font != nil {
		font.SetRGB(0, 0, 255)
		me.text = font.Write("starfish 0.12!")
		font.Free()
	} else {
		fmt.Println("Could not load LiberationSans-Bold.ttf.")
		return false
	}
	if me.box == nil {
		fmt.Println("Could not load box.png.")
		return false
	}
	me.anim.LoadImageSize("box.png", 70, 70)
	me.anim.LoadImageSize("dots.png", 70, 70)
	return true
}

func (me *Drawer) Draw(c *gfx.Canvas) {
	//clear screen
	c.SetRGB(0, 0, 0)
	c.FillRect(0, 0, gfx.DisplayWidth(), gfx.DisplayHeight())

	c.SetRGBA(0, 0, 255, 255)
	c.FillRect(42, 42, 100, 100)

	//draw box if it's not nil
	if me.box != nil {
		c.DrawImage(me.box, 200, 200)
		c.SetRGBA(0, 0, 0, 100)
		c.FillRect(200, 200, 100, 100)
	}
	c.DrawText(me.text, 400, 400)

	//Note: viewports may be nested
	c.PushViewport(42, 42, 500, 500)
	{
		//draw a green rect in a viewport
		c.SetRGBA(0, 255, 100, 127)
		c.FillRect(42, 42, 100, 100)
		c.SetRGB(0, 0, 0)
		c.FillRect(350, 200, 70, 70)
		c.DrawAnimation(me.anim, 350, 200)
	}
	c.PopViewport()
}

func main() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//For a fullscreen at your screens native resolution, simply use this line instead:
	//if !gfx.OpenDisplay(0, 0, true) {
	if !gfx.OpenDisplay(800, 600, false) {
		return
	}
	gfx.SetDisplayTitle("starfish example")

	var pane Drawer
	if !pane.init() {
		return
	}
	gfx.AddDrawer(&pane)
	quit := func() {
		gfx.CloseDisplay()
		pane.box.Free()
		pane.text.Free()
	}
	input.AddQuitFunc(quit)
	input.AddMouseWheelFunc(func(e input.MouseWheelEvent) {
		if e.Up {
			fmt.Println("Mouse wheel scrolling up.")
		} else {
			fmt.Println("Mouse wheel scrolling down.")
		}
	})
	input.AddMousePressFunc(func(e input.MouseEvent) {
		fmt.Println("Mouse Press!")
	})
	input.AddMouseReleaseFunc(func(e input.MouseEvent) {
		fmt.Println("Mouse Release!")
	})

	input.AddKeyPressFunc(func(e input.KeyEvent) {
		fmt.Println("Key Press!")
		if e.Key == input.Key_Escape {
			quit()
		}
	})
	input.AddKeyReleaseFunc(func(i input.KeyEvent) {
		fmt.Println("Key Release!")
	})

	gfx.Main()
}
