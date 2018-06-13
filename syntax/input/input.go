/* This package defines the various user functions related to input devices such as keyboard, mouse and joystick.
 * 
 * API:
 * 
 * // Returns a list of strings representing all of the keys that have been pressed since the last update.
 * // This is safe to add to an input string eg. If shift is down, alphabetical strings will be capitalised.
 * keys.pressed() 
 * ['d', len(#keys), for key in keys { len(key), bytes(key)... } ]
 */

package input

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/list"
import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/syntax/subjective"

var Keys = compiler.Expression{
	Name: compiler.Translatable {
		compiler.English: "keys",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		
		c.Expecting(".")
		c.Expecting("pressed")
		c.Expecting("(")
		c.Expecting(")")

		c.Call(&Keys_Pressed)
		
		return list.Of(text.Type)
	},
}

var Keys_Pressed = compiler.Function {
	Name: compiler.NoTranslation("keys_pressed"),
	
	Flags: []compiler.Flag{ subjective.Flag },

	Compile: func(c *compiler.Compiler) {
		c.List()
		
		c.SwapOutput()
		
		c.Native("go", `
var Keyboard_Buffer [256]bool
var Last_Keyboard_Buffer [256]bool
func Update_Keyboard_Buffer() {
	Keyboard_Buffer['1'] = ebiten.IsKeyPressed(ebiten.Key1)
	Keyboard_Buffer['2'] = ebiten.IsKeyPressed(ebiten.Key2)
	Keyboard_Buffer['3'] = ebiten.IsKeyPressed(ebiten.Key3)
	Keyboard_Buffer['4'] = ebiten.IsKeyPressed(ebiten.Key4)
	Keyboard_Buffer['5'] = ebiten.IsKeyPressed(ebiten.Key5)
	Keyboard_Buffer['6'] = ebiten.IsKeyPressed(ebiten.Key6)
	Keyboard_Buffer['7'] = ebiten.IsKeyPressed(ebiten.Key7)
	Keyboard_Buffer['8'] = ebiten.IsKeyPressed(ebiten.Key8) 
	Keyboard_Buffer['9'] = ebiten.IsKeyPressed(ebiten.Key9) 
	Keyboard_Buffer['a'] = ebiten.IsKeyPressed(ebiten.KeyA) 
	Keyboard_Buffer['b'] = ebiten.IsKeyPressed(ebiten.KeyB) 
	Keyboard_Buffer['c'] = ebiten.IsKeyPressed(ebiten.KeyC) 
	Keyboard_Buffer['d'] = ebiten.IsKeyPressed(ebiten.KeyD) 
	Keyboard_Buffer['e'] = ebiten.IsKeyPressed(ebiten.KeyE) 
	Keyboard_Buffer['f'] = ebiten.IsKeyPressed(ebiten.KeyF) 
	Keyboard_Buffer['g'] = ebiten.IsKeyPressed(ebiten.KeyG) 
	Keyboard_Buffer['h'] = ebiten.IsKeyPressed(ebiten.KeyH) 
	Keyboard_Buffer['i'] = ebiten.IsKeyPressed(ebiten.KeyI) 
	Keyboard_Buffer['j'] = ebiten.IsKeyPressed(ebiten.KeyJ) 
	Keyboard_Buffer['k'] = ebiten.IsKeyPressed(ebiten.KeyK) 
	Keyboard_Buffer['l'] = ebiten.IsKeyPressed(ebiten.KeyL) 
	Keyboard_Buffer['m'] = ebiten.IsKeyPressed(ebiten.KeyM) 
	Keyboard_Buffer['n'] = ebiten.IsKeyPressed(ebiten.KeyN) 
	Keyboard_Buffer['o'] = ebiten.IsKeyPressed(ebiten.KeyO) 
	Keyboard_Buffer['p'] = ebiten.IsKeyPressed(ebiten.KeyP) 
	Keyboard_Buffer['q'] = ebiten.IsKeyPressed(ebiten.KeyQ) 
	Keyboard_Buffer['r'] = ebiten.IsKeyPressed(ebiten.KeyR) 
	Keyboard_Buffer['s'] = ebiten.IsKeyPressed(ebiten.KeyS) 
	Keyboard_Buffer['t'] = ebiten.IsKeyPressed(ebiten.KeyT) 
	Keyboard_Buffer['u'] = ebiten.IsKeyPressed(ebiten.KeyU) 
	Keyboard_Buffer['v'] = ebiten.IsKeyPressed(ebiten.KeyV) 
	Keyboard_Buffer['w'] = ebiten.IsKeyPressed(ebiten.KeyW) 
	Keyboard_Buffer['x'] = ebiten.IsKeyPressed(ebiten.KeyX) 
	Keyboard_Buffer['y'] = ebiten.IsKeyPressed(ebiten.KeyY) 
	Keyboard_Buffer['z'] = ebiten.IsKeyPressed(ebiten.KeyZ) 
	/*Keyboard_Buffer[ebiten.KeyAlt] = ebiten.IsKeyPressed(ebiten.KeyAlt) 
	Keyboard_Buffer[ebiten.KeyApostrophe] = ebiten.IsKeyPressed(ebiten.KeyApostrophe) 
	Keyboard_Buffer[ebiten.KeyBackslash] = ebiten.IsKeyPressed(ebiten.KeyBackslash) 
	Keyboard_Buffer[ebiten.KeyBackspace] = ebiten.IsKeyPressed(ebiten.KeyBackspace) 
	Keyboard_Buffer[ebiten.KeyCapsLock] = ebiten.IsKeyPressed(ebiten.KeyCapsLock) 
	Keyboard_Buffer[ebiten.KeyComma] = ebiten.IsKeyPressed(ebiten.KeyComma) 
	Keyboard_Buffer[ebiten.KeyControl] = ebiten.IsKeyPressed(ebiten.KeyControl) 
	Keyboard_Buffer[ebiten.KeyDelete] = ebiten.IsKeyPressed(ebiten.KeyDelete) 
	Keyboard_Buffer[ebiten.KeyDown] = ebiten.IsKeyPressed(ebiten.KeyDown) 
	Keyboard_Buffer[ebiten.KeyEnd] = ebiten.IsKeyPressed(ebiten.KeyEnd) 
	Keyboard_Buffer[ebiten.KeyEnter] = ebiten.IsKeyPressed(ebiten.KeyEnter) 
	Keyboard_Buffer[ebiten.KeyEqual] = ebiten.IsKeyPressed(ebiten.KeyEqual) 
	Keyboard_Buffer[ebiten.KeyEscape] = ebiten.IsKeyPressed(ebiten.KeyEscape) 
	Keyboard_Buffer[ebiten.KeyF1] = ebiten.IsKeyPressed(ebiten.KeyF1) 
	Keyboard_Buffer[ebiten.KeyF2] = ebiten.IsKeyPressed(ebiten.KeyF2) 
	Keyboard_Buffer[ebiten.KeyF3] = ebiten.IsKeyPressed(ebiten.KeyF3) 
	Keyboard_Buffer[ebiten.KeyF4] = ebiten.IsKeyPressed(ebiten.KeyF4) 
	Keyboard_Buffer[ebiten.KeyF5] = ebiten.IsKeyPressed(ebiten.KeyF5) 
	Keyboard_Buffer[ebiten.KeyF6] = ebiten.IsKeyPressed(ebiten.KeyF6) 
	Keyboard_Buffer[ebiten.KeyF7] = ebiten.IsKeyPressed(ebiten.KeyF7) 
	Keyboard_Buffer[ebiten.KeyF8] = ebiten.IsKeyPressed(ebiten.KeyF8) 
	Keyboard_Buffer[ebiten.KeyF9] = ebiten.IsKeyPressed(ebiten.KeyF9) 
	Keyboard_Buffer[ebiten.KeyF10] = ebiten.IsKeyPressed(ebiten.KeyF10) 
	Keyboard_Buffer[ebiten.KeyF11] = ebiten.IsKeyPressed(ebiten.KeyF11) 
	Keyboard_Buffer[ebiten.KeyF12] = ebiten.IsKeyPressed(ebiten.KeyF12) 
	Keyboard_Buffer[ebiten.KeyGraveAccent] = ebiten.IsKeyPressed(ebiten.KeyGraveAccent) 
	Keyboard_Buffer[ebiten.KeyHome] = ebiten.IsKeyPressed(ebiten.KeyHome) 
	Keyboard_Buffer[ebiten.KeyInsert] = ebiten.IsKeyPressed(ebiten.KeyInsert) 
	Keyboard_Buffer[ebiten.KeyKP0] = ebiten.IsKeyPressed(ebiten.KeyKP0) 
	Keyboard_Buffer[ebiten.KeyKP1] = ebiten.IsKeyPressed(ebiten.KeyKP1) 
	Keyboard_Buffer[ebiten.KeyKP2] = ebiten.IsKeyPressed(ebiten.KeyKP2) 
	Keyboard_Buffer[ebiten.KeyKP3] = ebiten.IsKeyPressed(ebiten.KeyKP3) 
	Keyboard_Buffer[ebiten.KeyKP4] = ebiten.IsKeyPressed(ebiten.KeyKP4) 
	Keyboard_Buffer[ebiten.KeyKP5] = ebiten.IsKeyPressed(ebiten.KeyKP5) 
	Keyboard_Buffer[ebiten.KeyKP6] = ebiten.IsKeyPressed(ebiten.KeyKP6) 
	Keyboard_Buffer[ebiten.KeyKP7] = ebiten.IsKeyPressed(ebiten.KeyKP7) 
	Keyboard_Buffer[ebiten.KeyKP8] = ebiten.IsKeyPressed(ebiten.KeyKP8) 
	Keyboard_Buffer[ebiten.KeyKP9] = ebiten.IsKeyPressed(ebiten.KeyKP9) 
	Keyboard_Buffer[ebiten.KeyKPAdd] = ebiten.IsKeyPressed(ebiten.KeyKPAdd) 
	Keyboard_Buffer[ebiten.KeyKPDecimal] = ebiten.IsKeyPressed(ebiten.KeyKPDecimal) 
	Keyboard_Buffer[ebiten.KeyKPDivide] = ebiten.IsKeyPressed(ebiten.KeyKPDivide) 
	Keyboard_Buffer[ebiten.KeyKPEnter] = ebiten.IsKeyPressed(ebiten.KeyKPEnter) 
	Keyboard_Buffer[ebiten.KeyKPEqual] = ebiten.IsKeyPressed(ebiten.KeyKPEqual) 
	Keyboard_Buffer[ebiten.KeyKPMultiply] = ebiten.IsKeyPressed(ebiten.KeyKPMultiply) 
	Keyboard_Buffer[ebiten.KeyKPSubtract] = ebiten.IsKeyPressed(ebiten.KeyKPSubtract) 
	Keyboard_Buffer[ebiten.KeyLeft] = ebiten.IsKeyPressed(ebiten.KeyLeft) 
	Keyboard_Buffer[ebiten.KeyLeftBracket] = ebiten.IsKeyPressed(ebiten.KeyLeftBracket) 
	Keyboard_Buffer[ebiten.KeyMenu] = ebiten.IsKeyPressed(ebiten.KeyMenu) 
	Keyboard_Buffer[ebiten.KeyMinus] = ebiten.IsKeyPressed(ebiten.KeyMinus) 
	Keyboard_Buffer[ebiten.KeyNumLock] = ebiten.IsKeyPressed(ebiten.KeyNumLock) 
	Keyboard_Buffer[ebiten.KeyPageDown] = ebiten.IsKeyPressed(ebiten.KeyPageDown) 
	Keyboard_Buffer[ebiten.KeyPageUp] = ebiten.IsKeyPressed(ebiten.KeyPageUp) 
	Keyboard_Buffer[ebiten.KeyPause] = ebiten.IsKeyPressed(ebiten.KeyPause) 
	Keyboard_Buffer[ebiten.KeyPeriod] = ebiten.IsKeyPressed(ebiten.KeyPeriod) 
	Keyboard_Buffer[ebiten.KeyPrintScreen] = ebiten.IsKeyPressed(ebiten.KeyPrintScreen) 
	Keyboard_Buffer[ebiten.KeyRight] = ebiten.IsKeyPressed(ebiten.KeyRight) 
	Keyboard_Buffer[ebiten.KeyRightBracket] = ebiten.IsKeyPressed(ebiten.KeyRightBracket) 
	Keyboard_Buffer[ebiten.KeyScrollLock] = ebiten.IsKeyPressed(ebiten.KeyScrollLock) 
	Keyboard_Buffer[ebiten.KeySemicolon] = ebiten.IsKeyPressed(ebiten.KeySemicolon) 
	Keyboard_Buffer[ebiten.KeyShift] = ebiten.IsKeyPressed(ebiten.KeyShift) 
	Keyboard_Buffer[ebiten.KeySlash] = ebiten.IsKeyPressed(ebiten.KeySlash) 
	Keyboard_Buffer[ebiten.KeySpace] = ebiten.IsKeyPressed(ebiten.KeySpace) 
	Keyboard_Buffer[ebiten.KeyTab] = ebiten.IsKeyPressed(ebiten.KeyTab) 
	Keyboard_Buffer[ebiten.KeyUp] = ebiten.IsKeyPressed(ebiten.KeyUp) 
	Keyboard_Buffer[ebiten.KeyMax] = ebiten.IsKeyPressed(ebiten.KeyMax)*/
}

`)

c.Native("py", `
keyboard_buffer = {}
def on_key_press(symbol, modifiers):	
	global keyboard_buffer
	if (symbol >= 97 and symbol <= 122):
		keyboard_buffer[chr(symbol)] = True
	
grate_window.on_key_press = on_key_press
`)
		
		c.SwapOutput()

		c.Native("go", `
		for i := range Keyboard_Buffer {
			if Keyboard_Buffer[i] && !Last_Keyboard_Buffer[i] {
				runtime.Stack = append(runtime.Stack, Int{})
				runtime.Lists = append(runtime.Lists, &List{Mixed: []Int{Int{Small:int64(i)}}, Bytes: []byte{byte(i)}})
				runtime.HeapList()
				runtime.Lists[len(runtime.Lists)-1].Put(runtime.Stack[len(runtime.Stack)-1])
				runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]
			}
		}
`)
		c.Native("py", `
	for key in keyboard_buffer:
		if keyboard_buffer[key]:
			runtime.Lists.append([ord(key)])
			runtime.Stack.append(0)
			runtime.heaplist()
			runtime.Lists[-1].append(runtime.Stack.pop())
`)
	},
}
