/*
 * Copyright 2011-2012 starfish authors
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <SDL2/SDL_ttf.h>
#include "sdl.h"

SDL_Window *screen;

SDL_Window *openDisplay(int w, int h, int full) {
	SDL_Init(SDL_INIT_VIDEO);
	TTF_Init();
	unsigned flags = 0;
	if (full) {
		flags |= SDL_WINDOW_FULLSCREEN;
	}
	screen = SDL_CreateWindow("", SDL_WINDOWPOS_UNDEFINED, SDL_WINDOWPOS_UNDEFINED, w, h, flags);
	return screen;
}

void closeDisplay() {
	SDL_DestroyWindow(screen);
	screen = NULL;
	SDL_Quit();
}

int eventType(SDL_Event *e) {
	return e->type;
}

int eventKey(SDL_Event *e) {
	return e->key.keysym.sym;
}

int eventMouseButton(SDL_Event *e) {
	return e->button.button;
}

int eventMouseX(SDL_Event *e) {
	return e->button.x;
}

int eventMouseY(SDL_Event *e) {
	return e->button.y;
}
