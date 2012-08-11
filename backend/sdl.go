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
package backend

/*
#cgo LDFLAGS: -lSDL -lSDL_ttf -lSDL_image -lSDL_gfx
#include <SDL/SDL.h>
#include "sdl.h"
#include "SDL/SDL_rotozoom.h"
#include "SDL/SDL_gfxPrimitives.h"
#include "SDL/SDL_image.h"
#include "SDL/SDL_ttf.h"
*/
import "C"

var drawFunc func()
var screen *C.SDL_Surface

func OpenDisplay(w, h int, full bool) {
	var i C.int = 0
	if full {
		i = 1
	}
	screen = C.openDisplay(C.int(w), C.int(h), i)
}

func CloseDisplay() {
	C.closeDisplay()
}

func SetDrawFunc(f func()) {
	drawFunc = f
}

//Sets the title of the window.
func SetDisplayTitle(title string) {
	C.SDL_WM_SetCaption(C.CString(title), C.CString(""))
}

//Returns the width of the display window.
func DisplayWidth() int {
	return int(screen.w)
}

//Returns the height of the display window.
func DisplayHeight() int {
	return int(screen.h)
}

//Used to manually draw the screen.
func Draw() {
	drawFunc()
	C.SDL_Flip(screen)
}

//MISC

func sdl_Rect(x, y, width, height int) C.SDL_Rect {
	var r C.SDL_Rect
	r.x = C.Sint16(x)
	r.y = C.Sint16(y)
	r.w = C.Uint16(width)
	r.h = C.Uint16(height)
	return r
}

//An RGB color representation.
type Color struct {
	Red, Green, Blue, Alpha byte
}

func (me *Color) toSDL_Color() C.SDL_Color {
	return C.SDL_Color{C.Uint8(me.Red), C.Uint8(me.Green), C.Uint8(me.Blue), C.Uint8(me.Alpha)}
}

func (me *Color) toUint32() uint32 {
	return (uint32(me.Red) << 16) | (uint32(me.Green) << 8) | uint32(me.Blue)
}

//IMAGE HANDLING

type Image struct {
	surface *C.SDL_Surface
}

func (me *Image) W() int {
	return int(me.surface.w)
}

func (me *Image) H() int {
	return int(me.surface.h)
}

func LoadImage(path string) *Image {
	tmp := C.IMG_Load(C.CString(path))
	i := C.SDL_DisplayFormatAlpha(tmp)
	C.SDL_FreeSurface(tmp)

	out := new(Image)
	out.surface = i
	return out
}

func FreeImage(img *Image) {
	C.SDL_FreeSurface(img.surface)
}

func ResizeAngleOf(image *Image, angle float64, width, height int) *Image {
	img := image.surface
	if img.w == 0 || img.h == 0 {
		return nil
	}
	retval := new(Image)
	xstretch := C.double(float64(width) / float64(img.w))
	ystretch := C.double(float64(height) / float64(img.h))
	retval.surface = C.rotozoomSurfaceXY(img, C.double(angle), xstretch, ystretch, 1)
	return retval
}

//TEXT HANDLING

type Font struct {
	font *C.TTF_Font
}

func LoadFont(path string, size int) *Font {
	font := new(Font)
	font.font = C.TTF_OpenFont(C.CString(path), C.int(size))
	return font
}

func FreeFont(val *Font) {
	C.TTF_CloseFont(val.font)
}

func (me *Font) WriteTo(text string, t *Image, c Color) bool {
	t.surface = C.TTF_RenderText_Blended(me.font, C.CString(text), c.toSDL_Color())
	return t.surface != nil
}

//GFX HANDLING

//Pushes a viewport to limit the drawing space to the given bounds within the current drawing space.
func SetClipRect(x, y, w, h int) {
	r := sdl_Rect(x, y, w, h)
	C.SDL_SetClipRect(screen, &r)
}

func FillRoundedRect(x, y, width, height, radius int, c Color) {
	r := sdl_Rect(x, y, width, height)
	C.roundedBoxRGBA(screen, C.Sint16(r.x), C.Sint16(r.y), C.Sint16(int(r.x)+int(r.w)), C.Sint16(int(r.y)+int(r.h)), C.Sint16(radius), C.Uint8(c.Red), C.Uint8(c.Green), C.Uint8(c.Blue), C.Uint8(c.Alpha))
}

func FillRect(x, y, width, height int, c Color) {
	r := sdl_Rect(x, y, width, height)
	C.boxRGBA(screen, C.Sint16(r.x), C.Sint16(r.y), C.Sint16(int(r.x)+int(r.w)), C.Sint16(int(r.y)+int(r.h)), C.Uint8(c.Red), C.Uint8(c.Green), C.Uint8(c.Blue), C.Uint8(c.Alpha))
}

//Draws the image at the given coordinates.
func DrawImage(img *Image, x, y int) {
	C.SDL_SetAlpha(img.surface, C.SDL_SRCALPHA, 255)
	var dest C.SDL_Rect
	dest.x = C.Sint16(x)
	dest.y = C.Sint16(y)
	C.SDL_BlitSurface(img.surface, nil, screen, &dest)
}

//INPUT HANDLING

var running = true

func HandleInput() {
	running = true
	scrollFunc := func(b bool, x, y int) {
		var e MouseWheelEvent
		e.Up = b
		e.X = x
		e.Y = y
		MouseWheelFunc(e)
	}
	for running {
		var e C.SDL_Event
		C.SDL_WaitEvent(&e)
		switch C.eventType(&e) {
		case C.SDL_QUIT:
			go QuitFunc()
			running = false
		case C.SDL_KEYDOWN:
			go func() {
				var ke KeyEvent
				ke.Key = int(C.eventKey(&e))
				go KeyDown(ke)
			}()
		case C.SDL_KEYUP:
			go func() {
				var ke KeyEvent
				ke.Key = int(C.eventKey(&e))
				go KeyUp(ke)
			}()
		case C.SDL_MOUSEBUTTONDOWN:
			x := int(C.eventMouseX(&e))
			y := int(C.eventMouseY(&e))
			switch C.eventMouseButton(&e) {
			case C.SDL_BUTTON_WHEELUP:
				go scrollFunc(true, x, y)
			case C.SDL_BUTTON_WHEELDOWN:
				go scrollFunc(false, x, y)
			default:
				go func() {
					var me MouseEvent
					me.X = x
					me.Y = y
					me.Button = int(C.eventMouseButton(&e))
					go MouseButtonDown(me)
				}()
			}
		case C.SDL_MOUSEBUTTONUP:
			x := int(C.eventMouseX(&e))
			y := int(C.eventMouseY(&e))
			switch C.eventMouseButton(&e) {
			case C.SDL_BUTTON_WHEELUP:
				go scrollFunc(true, x, y)
			case C.SDL_BUTTON_WHEELDOWN:
				go scrollFunc(false, x, y)
			default:
				go func() {
					var me MouseEvent
					me.Button = int(C.eventMouseButton(&e))
					me.X = int(C.eventMouseX(&e))
					me.Y = int(C.eventMouseY(&e))
					me.Button = int(C.eventMouseButton(&e))
					go MouseButtonUp(me)
				}()
			}
		}
	}
}

//KEY DEFINITIONS

const (
	Key_a int = C.SDLK_a
	Key_b     = C.SDLK_b
	Key_c     = C.SDLK_c
	Key_d     = C.SDLK_d
	Key_e     = C.SDLK_e
	Key_f     = C.SDLK_f
	Key_g     = C.SDLK_g
	Key_h     = C.SDLK_h
	Key_i     = C.SDLK_i
	Key_j     = C.SDLK_j
	Key_k     = C.SDLK_k
	Key_l     = C.SDLK_l
	Key_m     = C.SDLK_m
	Key_n     = C.SDLK_n
	Key_o     = C.SDLK_o
	Key_p     = C.SDLK_p
	Key_q     = C.SDLK_q
	Key_r     = C.SDLK_r
	Key_s     = C.SDLK_s
	Key_t     = C.SDLK_t
	Key_u     = C.SDLK_u
	Key_v     = C.SDLK_v
	Key_w     = C.SDLK_w
	Key_x     = C.SDLK_x
	Key_y     = C.SDLK_y
	Key_z     = C.SDLK_z

	Key_0 = C.SDLK_0
	Key_1 = C.SDLK_1
	Key_2 = C.SDLK_2
	Key_3 = C.SDLK_3
	Key_4 = C.SDLK_4
	Key_5 = C.SDLK_5
	Key_6 = C.SDLK_6
	Key_7 = C.SDLK_7
	Key_8 = C.SDLK_8
	Key_9 = C.SDLK_9

	Key_Colon        = C.SDLK_COLON
	Key_SemiColon    = C.SDLK_SEMICOLON
	Key_LessThan     = C.SDLK_LESS
	Key_Equals       = C.SDLK_EQUALS
	Key_GreaterThan  = C.SDLK_GREATER
	Key_QuestionMark = C.SDLK_QUESTION
	Key_At           = C.SDLK_AT
	Key_LeftBracket  = C.SDLK_LEFTBRACKET
	Key_RightBracket = C.SDLK_RIGHTBRACKET
	Key_Caret        = C.SDLK_CARET
	Key_Underscore   = C.SDLK_UNDERSCORE
	Key_BackQuote    = C.SDLK_BACKQUOTE
	Key_Backspace    = C.SDLK_BACKSPACE
	Key_Tab          = C.SDLK_TAB
	Key_Clear        = C.SDLK_CLEAR
	Key_Enter        = C.SDLK_RETURN
	Key_Pause        = C.SDLK_PAUSE
	Key_Escape       = C.SDLK_ESCAPE
	Key_Space        = C.SDLK_SPACE
	Key_ExclaimMark  = C.SDLK_EXCLAIM
	Key_DoubleQuote  = C.SDLK_QUOTEDBL
	Key_Hash         = C.SDLK_HASH
	Key_Dollar       = C.SDLK_DOLLAR
	Key_LeftParen    = C.SDLK_LEFTPAREN
	Key_RightParen   = C.SDLK_RIGHTPAREN
	Key_Asterisk     = C.SDLK_ASTERISK
	Key_Plus         = C.SDLK_PLUS
	Key_Comma        = C.SDLK_COMMA
	Key_Minus        = C.SDLK_MINUS
	Key_Period       = C.SDLK_PERIOD
	Key_Slash        = C.SDLK_SLASH
	Key_Delete       = C.SDLK_DELETE
	Key_Up           = C.SDLK_UP
	Key_Down         = C.SDLK_DOWN
	Key_Left         = C.SDLK_LEFT
	Key_Right        = C.SDLK_RIGHT
	Key_RCtrl        = C.SDLK_RCTRL
	Key_LCtrl        = C.SDLK_LCTRL
)
