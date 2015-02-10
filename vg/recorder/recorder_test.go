// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recorder

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gonum/plot/vg"
)

func TestRecorder(t *testing.T) {
	r := New(72, []Action{&FillString{Font: "Foo", Size: 12, X: 0, Y: 10, String: "Bar"}})
	r.Scale(1, 2)
	r.Rotate(0.72)
	r.KeepCaller = true
	r.Stroke(vg.Path{{Type: vg.MoveComp, X: 3, Y: 4}})
	for i, a := range r.Actions {
		if i < len(r.Base) {
			if !reflect.DeepEqual(a, r.Base[i]) {
				t.Errorf("unexpected mismatch between base and actions:\n\tgot: %#v\n\twant: %#v", a, r.Base[i])
			}
		}
		if got := a.VgCall(); !strings.HasSuffix(got, want[i]) {
			t.Errorf("unexpected actions:\n\tgot: %#v\n\twant: %#v", got, want[i])
		}
	}
}

var want = []string{
	`FillString("Foo", 12, 0, 10, "Bar")`,
	`Scale(1, 2)`,
	`Rotate(0.72)`,
	`github.com/gonum/plot/vg/recorder/recorder_test.go:20 Stroke(vg.Path{vg.PathComp{Type:0, X:3, Y:4, Radius:0, Start:0, Angle:0}})`,
}
