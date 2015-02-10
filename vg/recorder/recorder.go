// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package recorder provides support for vector graphics serialization
// it is intended to be used as a testing helper.
package recorder

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/gonum/plot/vg"
)

var _ vg.Canvas = (*Canvas)(nil)

// Canvas implements vg.Canvas operation serialization.
type Canvas struct {
	// Resultion holds the canvas resolution in DPI.
	Resolution float64

	// Base holds the starting state of the canvas.
	Base []Action

	// Actions holds the complete state of the canvas.
	Actions []Action

	// KeepCaller indicates whether the Canvas will
	// retain caller information for the actions.
	KeepCaller bool
}

// Action is a vector graphics action as defined by the
// vg.Canvas interface. Each method of vg.Canvas must have
// a corresponding Action type.
type Action interface {
	VgCall() string
	callRecorder() *callRecorder
}

type callRecorder struct {
	haveCaller bool
	file       string
	line       int
}

func (c *callRecorder) caller() {
	_, c.file, c.line, c.haveCaller = runtime.Caller(3)
}

func (c callRecorder) String() string {
	if !c.haveCaller {
		return ""
	}
	return fmt.Sprintf("%s:%d ", c.file, c.line)
}

// New returns a new Canvas pre-filled with the given base actions.
func New(dpi float64, base []Action) *Canvas {
	return &Canvas{
		Resolution: dpi,
		Actions:    base,
		Base:       base,
	}
}

// Reset resets the Canvas to the base state.
func (c *Canvas) Reset() {
	c.Actions = c.Actions[:len(c.Base)]
	copy(c.Actions, c.Base)
}

func (c *Canvas) append(a Action) {
	if c.KeepCaller {
		a.callRecorder().caller()
	}
	c.Actions = append(c.Actions, a)
}

// SetWidth corresponds to the vg.Canvas.SetWidth method.
type SetWidth struct {
	vg.Length

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *SetWidth) VgCall() string {
	return fmt.Sprintf("%sSetWidth(%v)", a.cr, a.Length)
}

func (a *SetWidth) callRecorder() *callRecorder {
	return &a.cr
}

// SetLineWidth implements the SetLineWidth method of the vg.Canvas interface.
func (c *Canvas) SetLineWidth(w vg.Length) {
	c.append(&SetWidth{Length: w})
}

// SetLineDash corresponds to the vg.Canvas.SetLineDash method.
type SetLineDash struct {
	Dashes  []vg.Length
	Offsets vg.Length

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *SetLineDash) VgCall() string {
	return fmt.Sprintf("%sSetLineDash(%#v, %v)", a.cr, a.Dashes, a.Offsets)
}

func (a *SetLineDash) callRecorder() *callRecorder {
	return &a.cr
}

// SetLineDash implements the SetLineDash method of the vg.Canvas interface.
func (c *Canvas) SetLineDash(dashes []vg.Length, offs vg.Length) {
	c.append(&SetLineDash{
		Dashes:  append([]vg.Length(nil), dashes...),
		Offsets: offs,
	})
}

// SetColor corresponds to the vg.Canvas.SetColor method.
type SetColor struct {
	color.Color

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *SetColor) VgCall() string {
	return fmt.Sprintf("%sSetColor(%#v)", a.cr, a.Color)
}

func (a *SetColor) callRecorder() *callRecorder {
	return &a.cr
}

// SetColor implements the SetColor method of the vg.Canvas interface.
func (c *Canvas) SetColor(col color.Color) {
	c.append(&SetColor{Color: col})
}

// Rotate corresponds to the vg.Canvas.Rotate method.
type Rotate struct {
	Angle float64

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *Rotate) VgCall() string {
	return fmt.Sprintf("%sRotate(%v)", a.cr, a.Angle)
}

func (a *Rotate) callRecorder() *callRecorder {
	return &a.cr
}

// Rotate implements the Rotate method of the vg.Canvas interface.
func (c *Canvas) Rotate(a float64) {
	c.append(&Rotate{Angle: a})
}

// Translate corresponds to the vg.Canvas.Translate method.
type Translate struct {
	X, Y vg.Length

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *Translate) VgCall() string {
	return fmt.Sprintf("%sTranslate(%v, %v)", a.cr, a.X, a.Y)
}

func (a *Translate) callRecorder() *callRecorder {
	return &a.cr
}

// Translate implements the Translate method of the vg.Canvas interface.
func (c *Canvas) Translate(x, y vg.Length) {
	c.append(&Translate{X: x, Y: y})
}

// Scale corresponds to the vg.Canvas.Scale method.
type Scale struct {
	X, Y float64

	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *Scale) VgCall() string {
	return fmt.Sprintf("%sScale(%v, %v)", a.cr, a.X, a.Y)
}

func (a *Scale) callRecorder() *callRecorder {
	return &a.cr
}

// Scale implements the Scale method of the vg.Canvas interface.
func (c *Canvas) Scale(x, y float64) {
	c.append(&Scale{X: x, Y: y})
}

// Push corresponds to the vg.Canvas.Push method.
type Push struct {
	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *Push) VgCall() string {
	return fmt.Sprintf("%sPush()", a)
}

func (a *Push) callRecorder() *callRecorder {
	return &a.cr
}

// Push implements the Push method of the vg.Canvas interface.
func (c *Canvas) Push() {
	c.append(&Push{})
}

// Pop corresponds to the vg.Canvas.Pop method.
type Pop struct {
	cr callRecorder
}

func (a *Pop) callRecorder() *callRecorder {
	return &a.cr
}

// VgCall returns the method call that generated the action.
func (a *Pop) VgCall() string {
	return fmt.Sprintf("%sPop()", a)
}

// Pop implements the Pop method of the vg.Canvas interface.
func (c *Canvas) Pop() {
	c.append(&Pop{})
}

// Stroke corresponds to the vg.Canvas.Stroke method.
type Stroke struct {
	vg.Path

	cr callRecorder
}

func (a *Stroke) callRecorder() *callRecorder {
	return &a.cr
}

// VgCall returns the method call that generated the action.
func (a *Stroke) VgCall() string {
	return fmt.Sprintf("%sStroke(%#v)", a.cr, a.Path)
}

// Stroke implements the Stroke method of the vg.Canvas interface.
func (c *Canvas) Stroke(path vg.Path) {
	c.append(&Stroke{Path: append(vg.Path(nil), path...)})
}

// Fill corresponds to the vg.Canvas.Fill method.
type Fill struct {
	vg.Path

	cr callRecorder
}

func (a *Fill) callRecorder() *callRecorder {
	return &a.cr
}

// VgCall returns the method call that generated the action.
func (a *Fill) VgCall() string {
	return fmt.Sprintf("%sFill(%#v)", a.cr, a.Path)
}

// Fill implements the Fill method of the vg.Canvas interface.
func (c *Canvas) Fill(path vg.Path) {
	c.append(&Fill{Path: append(vg.Path(nil), path...)})
}

// FillString corresponds to the vg.Canvas.FillString method.
type FillString struct {
	Font   string
	Size   vg.Length
	X, Y   vg.Length
	String string

	cr callRecorder
}

func (a *FillString) callRecorder() *callRecorder {
	return &a.cr
}

// VgCall returns the method call that generated the action.
func (a *FillString) VgCall() string {
	return fmt.Sprintf("%sFillString(%q, %v, %v, %v, %q)", a.cr, a.Font, a.Size, a.X, a.Y, a.String)
}

// FillString implements the FillString method of the vg.Canvas interface.
func (c *Canvas) FillString(font vg.Font, x, y vg.Length, str string) {
	c.append(&FillString{
		Font: font.Name(),
		Size: font.Size,
		X:    x, Y: y,
		String: str,
	})
}

// DPI corresponds to the vg.Canvas.DPI method.
type DPI struct {
	cr callRecorder
}

// VgCall returns the method call that generated the action.
func (a *DPI) VgCall() string {
	return fmt.Sprintf("%sDPI()", a)
}

func (a *DPI) callRecorder() *callRecorder {
	return &a.cr
}

// DPI implements the DPI method of the vg.Canvas interface.
func (c *Canvas) DPI() float64 {
	c.append(&DPI{})
	return c.Resolution
}
