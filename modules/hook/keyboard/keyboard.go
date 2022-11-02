package keyboard

import "github.com/qxdn/keyboard-wordcloud/modules/hook/types"

type HookHandler func(c chan<- types.KeyboardEvent) types.HOOKPROC

func Install(ch chan<- types.KeyboardEvent) error {
	return install(ch)
}

func Uninstall() error {
	return uninstall()
}
