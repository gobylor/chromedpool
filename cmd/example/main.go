package main

import (
	"log"
	"sync"

	"github.com/gobylor/chromedpool"
)

func main() {
	pool, err := chromedpool.NewTabPool(
		chromedpool.WithMaxTabs(3),
		// Setting false to visualize the run, defaults to true.
		chromedpool.WithHeadless(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	urlPrefix := "https://en.wikipedia.org/wiki/"
	keywords := []string{"china", "japan", "america", "germany", "france", "italy", "spain", "russia", "india", "brazil"}

	var wg sync.WaitGroup
	for _, k := range keywords {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			tab := pool.GetTab()
			defer pool.PutTab(tab)
			if err := tab.Navigate(url); err != nil {
				log.Printf("Error navigating to %s: %v", url, err)
				return
			}
		}(urlPrefix + k)
	}

	wg.Wait()
}
