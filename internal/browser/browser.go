package browser

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/Lunarisnia/simple-browser/internal/url"
)

type Browser struct {
	app          fyne.App
	mainWindow   fyne.Window
	displayList  []*DisplayObject
	drawnContent []DrawnObject
	scrollbar    fyne.CanvasObject

	Width  int
	Height int
	Scroll float32
}

type DrawnObject struct {
	Type      string
	originalX float32
	originalY float32
	object    fyne.CanvasObject
}

func New(width int, height int) *Browser {
	a := app.New()
	mainWindow := a.NewWindow("Ignis")
	mainWindow.SetMaster()

	mainWindow.Resize(fyne.NewSize(float32(width),
		float32(height)))
	mainWindow.CenterOnScreen()
	return &Browser{
		app:          a,
		Scroll:       0.0,
		mainWindow:   mainWindow,
		drawnContent: make([]DrawnObject, 0),
		displayList:  make([]*DisplayObject, 0),
		Width:        width,
		Height:       height,
	}
}

func (b *Browser) Run() {
	b.mainWindow.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyDown {
			lastObject := b.drawnContent[len(b.drawnContent)-1]
			if lastObject.object.Position().Y <= float32(b.Height) {
				return
			}
			b.Scroll += 18.0
			b.Scroll = min(b.Scroll, lastObject.originalY)

			maxScrollPosition := lastObject.originalY - float32(b.Height)
			b.scrollbar.Move(fyne.NewPos(b.scrollbar.Position().X, b.Scroll/maxScrollPosition*(float32(b.Height)-b.scrollbar.Size().Height)))
			if b.scrollbar.Position().Y >= float32(b.Height)-b.scrollbar.Size().Height {
				b.scrollbar.Move(fyne.NewPos(b.scrollbar.Position().X, float32(b.Height)-b.scrollbar.Size().Height))
			}

			for _, drawn := range b.drawnContent {
				newPos := fyne.NewPos(drawn.originalX, drawn.originalY-b.Scroll)
				drawn.object.Move(newPos)
			}
		}
		if ke.Name == fyne.KeyUp {
			firstObject := b.drawnContent[0].object
			if firstObject.Position().Y >= float32(b.Height)/8.0 {
				return
			}
			b.Scroll -= 19.0
			b.Scroll = max(b.Scroll, 0.0)

			lastObject := b.drawnContent[len(b.drawnContent)-1]
			maxScrollPosition := lastObject.originalY - float32(b.Height)
			b.scrollbar.Move(fyne.NewPos(b.scrollbar.Position().X, b.Scroll/maxScrollPosition*(float32(b.Height)-b.scrollbar.Size().Height)))
			if b.scrollbar.Position().Y >= float32(b.Height)-b.scrollbar.Size().Height {
				b.scrollbar.Move(fyne.NewPos(b.scrollbar.Position().X, float32(b.Height)-b.scrollbar.Size().Height))
			}

			for _, drawn := range b.drawnContent {
				newPos := fyne.NewPos(drawn.originalX, drawn.originalY-b.Scroll)
				drawn.object.Move(newPos)
			}
		}
	})
	b.mainWindow.ShowAndRun()
}

func (b *Browser) Load(path url.URL) {
	body, err := url.Load(path)
	if err != nil {
		// TODO: Return proper error
		log.Fatal(err)
	}

	b.displayList = b.layout(body)
	b.Draw()
}

func (b *Browser) Draw() {
	content := container.NewWithoutLayout()

	for _, d := range b.displayList {
		text := canvas.NewText(d.Char, color.White)
		pos := fyne.NewPos(d.X, d.Y-b.Scroll)
		text.Move(pos)
		content.Add(text)
		b.drawnContent = append(b.drawnContent, DrawnObject{
			originalX: pos.X,
			originalY: pos.Y,
			object:    text,
		})
	}

	scrollbar := canvas.NewRectangle(color.White)
	scrollbar.Resize(fyne.NewSize(100.0, 100.0))
	scrollbar.Move(fyne.NewPos(float32(b.Width)-105.0, 0.0))
	b.scrollbar = scrollbar

	content.Add(scrollbar)

	b.mainWindow.SetContent(content)
}

type DisplayObject struct {
	Char string
	X    float32
	Y    float32
}

func (b *Browser) layout(body string) []*DisplayObject {
	displayList := make([]*DisplayObject, 0)

	xStep, yStep := float32(18.0), float32(19.0)
	cursorX, cursorY := xStep, yStep
	for _, r := range body {
		if r == '\n' {
			cursorY += yStep * 2.0
		}
		text := DisplayObject{
			Char: string(r),
			X:    cursorX,
			Y:    cursorY,
		}
		displayList = append(displayList, &text)
		cursorX += xStep
		lastCharacterOnRow := float32(b.Width) - xStep
		if cursorX >= lastCharacterOnRow-100.0 {
			cursorY += yStep
			cursorX = xStep
		}
	}

	return displayList
}
