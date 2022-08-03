package browser

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

type Browser struct {
	Context context.Context
	Close   context.CancelFunc
}

func InitializeBrowser() *Browser {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	return &Browser{
		Context: ctx,
		Close:   cancel,
	}
}

func (b *Browser) Screenshot(url, filename, element string, buffer *[]byte) {
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Screenshot(element, buffer, chromedp.NodeVisible),
	}

	err := chromedp.Run(b.Context, tasks)
	if err != nil {
		log.Fatalf("[Browser] Unable to screenshot element at %s: %+v\n", url, err)
	}

	err = ioutil.WriteFile(filename, *buffer, 0o644)
	if err != nil {
		log.Fatalf("[Browser] Unable to save screenshot to local file: %+v\n", err)
	}
}
