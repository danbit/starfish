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
package graphics

type resourceNode struct {
	uses int
	rsrc interface{}
}

type resourceKey interface {
	String() string
}

type resourceCatalog struct {
	rsrcs  map[string]*resourceNode
	out    chan interface{}
	in     chan interface{}
	load   func(resourceKey) (interface{}, bool)
	delete func(resourceKey, interface{})
}

func newResourceCatalog(load func(resourceKey) (interface{}, bool), delete func(resourceKey, interface{})) (r resourceCatalog) {
	r.rsrcs = make(map[string]*resourceNode)
	r.out = make(chan interface{})
	r.in = make(chan interface{})
	r.load = load
	r.delete = delete
	go r.run()
	return r
}

func (me *resourceCatalog) checkout(key resourceKey) interface{} {
	me.out <- key
	return <-me.out
}

func (me *resourceCatalog) checkin(key resourceKey) {
	me.in <- key
	<-me.in
}

func (me *resourceCatalog) run() {
	for {
		select {
		case input := <-me.out: //checkout
			key := input.(resourceKey)
			i, ok := me.rsrcs[key.String()]
			if ok {
				i.uses++
				me.out <- i.rsrc
			} else {
				tmp, ok := me.load(key)
				if ok {
					i = new(resourceNode)
					i.rsrc = tmp
					i.uses++
					me.rsrcs[key.String()] = i
					me.out <- i.rsrc
				} else {
					me.out <- nil
				}
			}
		case input := <-me.in: //checkin
			key := input.(resourceKey)
			i, ok := me.rsrcs[key.String()]
			if ok {
				i.uses--
				if i.uses == 0 {
					me.delete(key, me.rsrcs[key.String()].rsrc)
					me.rsrcs[key.String()] = nil, false
				}
				me.in <- true
			} else {
				me.in <- false
			}
		}
	}
}
