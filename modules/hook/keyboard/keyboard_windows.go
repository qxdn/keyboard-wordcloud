//go:build windows

package keyboard

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/qxdn/keyboard-wordcloud/modules/hook/types"
	"github.com/qxdn/keyboard-wordcloud/modules/hook/win32"
)

var hHook struct {
	sync.Mutex
	Pointer uintptr
}

// wparam,lparam detail: https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)
func DefaultHookHandler(c chan<- types.KeyboardEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			c <- types.KeyboardEvent{
				Message:         types.Message(wParam),
				KBDLLHOOKSTRUCT: *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
		}

		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}

func install(ch chan<- types.KeyboardEvent) error {
	hHook.Lock()
	defer hHook.Unlock()
	if hHook.Pointer != 0 {
		return fmt.Errorf("keyboard: already installed keyboard hook")
	}
	if ch == nil {
		return fmt.Errorf("keyboard: chan must not nil")
	}
	fn := DefaultHookHandler(ch)

	go func() {
		hhook := win32.SetWindowsHookEx(
			types.WH_KEYBOARD_LL, // use low level to listen global
			syscall.NewCallback(fn),
			0,
			0)
		if hhook == 0 {
			// return value is null
			panic("keyboard: hook install fail")
		}
		hHook.Pointer = hhook

		var msg *types.MSG

		// This hook is called in the context of the thread that installed it. The call is made by sending a message to the thread that installed the hook. Therefore, the thread that installed the hook must have a message loop.
		// reference: https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)?redirectedfrom=MSDN
		for {
			if hHook.Pointer == 0 {
				break
			}
			if win32.GetMessage(&msg, 0, 0, 0) != 0 {
				win32.TranslateMessage(&msg)
				win32.DispatchMessage(&msg)
			}
		}

	}()

	return nil
}

func uninstall() error {
	hHook.Lock()
	defer hHook.Unlock()

	if !win32.UnhookWindowsHookEx(hHook.Pointer) {
		return fmt.Errorf("mouse: failed to uninstall hook function")
	}

	hHook.Pointer = 0

	return nil
}
