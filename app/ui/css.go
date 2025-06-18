package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const css = `
.error-input {
	background-color: #fff0f0;
	border: 2px solid #ff5555;
	border-radius: 3px;
	box-shadow: 0 0 5px rgba(255, 0, 0, 0.3);
}
.error-input:focus {
	border-color: #ff0000;
	box-shadow: 0 0 8px rgba(255, 0, 0, 0.5);
}
`

func (a *Application) initCSS() {
	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Fatal("Ошибка получения экрана:", err)
	}

	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Ошибка создания CSS провайдера:", err)
	}

	err = provider.LoadFromData(css)
	if err != nil {
		log.Fatal("Ошибка загрузки CSS:", err)
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

type styleProvider interface {
	GetStyleContext() (*gtk.StyleContext, error)
}

func HighlightInputField(widget styleProvider, isValid bool) {
	ctx, err := widget.GetStyleContext()
	if err != nil {
		log.Printf("Ошибка получения стиля: %v", err)
		return
	}

	if isValid {
		ctx.RemoveClass("error-input")
	} else {
		ctx.AddClass("error-input")
	}
}
