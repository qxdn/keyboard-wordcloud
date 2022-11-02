//go:build !windows

package keyboard

import (
	"fmt"
	"github.com/qxdn/keyboard-wordcloud/modules/hook/types"
)

func install(ch chan<- types.KeyboardEvent) error {
	return fmt.Errorf("keyboard: only support windows now")
}

func uninstall() error {
	return fmt.Errorf("keyboard: only support windows now")
}
