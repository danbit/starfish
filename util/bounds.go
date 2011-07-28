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
package util

/*
  Represents location and size attributes.
*/
type Bounds struct {
	Point
	Size
}

/*
  Sets this Bounds to the given coordinates and dimensions.
*/
func (me *Bounds) Set(x, y, width, height int) {
	me.X = x
	me.Y = y
	me.Width = width
	me.Height = height
}

/*
  Returns the x coordinate + the width.
*/
func (me *Bounds) X2() int {
	return me.X + me.Width
}

/*
  Returns the y coordinate + the height.
*/
func (me *Bounds) Y2() int {
	return me.Y + me.Height
}
