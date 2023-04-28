package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/manifoldco/promptui"
	"log"
	"strings"
	"time"
)

func main() {
	fmt.Println("Start of processing")
	// 上位の行をコメントアウトすることで下位行を有効化できます。
	// lower lines can be enabled by commenting out the higher lines.
	// level0: Create an instance of chrome
	//ctx, cancel := chromedp.NewContext(context.Background()) /*
	// level0-debug1: Instance created with logs
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf)) /*
	// level0-debug2: Create instances in no-headless mode
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("enable-automation", false),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//*/

	// level1: Access the page
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.Click(`//*[@id="gb"]/div/div[2]/a`, chromedp.NodeVisible),
	); err != nil {
		cancel()
		log.Fatal("err1-1: Failed login")
	}

	// level2: Enter credential
	var mailAddress string
	fmt.Printf("Enter your e-mail address or phone number: ")
	fmt.Scan(&mailAddress)
	password := passwdInputer("Enter your password")

	// level3: Processing ID
	if err := chromedp.Run(ctx,
		chromedp.Click(`//*[@id="identifierId"]`, chromedp.NodeVisible),
		input.InsertText(mailAddress),
		chromedp.Click(`//*[@id="identifierNext"]/div/button/div[3]`, chromedp.NodeVisible),
	); err != nil {
		cancel()
		log.Fatal("err3-1@login: Failed login")
	}

	// level3-2: Page Transition Confirmation
	time.Sleep(1 * time.Second)
	var url2CheckTransition string
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`body > div > div > div > div  > div`),
		chromedp.Location(&url2CheckTransition),
	); err != nil {
		cancel()
		log.Fatal("err3-2@login: Failed in page transition confirmation process")
	}
	if !strings.Contains(url2CheckTransition, "https://accounts.google.com/v3/signin/challenge/pwd?TL=") {
		cancel()
		log.Fatal("err3-3@login: Failed to load on email address input page")
	}

	// level4: Processing password
	if err := chromedp.Run(ctx,
		chromedp.Click(`//*[@id="password"]/div[1]/div/div[1]/input`, chromedp.NodeVisible),
		input.InsertText(password),
		chromedp.Click(`//*[@id="passwordNext"]/div/button/div[3]`, chromedp.NodeVisible),
	); err != nil {
		cancel()
		log.Fatal("err4-1@login: Failed to operate the login button")
	}

	// level4-2: Page Transition Confirmation
	time.Sleep(10 * time.Second)
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`body > div > div > div > div  > div`),
		chromedp.Location(&url2CheckTransition),
	); err != nil {
		cancel()
		log.Fatal("err4-2@login: Failed in page transition confirmation process")
	}
	fmt.Println(url2CheckTransition)
	if !strings.Contains(url2CheckTransition, "https://www.google.com/") {
		cancel()
		log.Fatal("err4-3@login: Failed to load on password input page")
	}

	fmt.Println("End of processing")
}

func passwdInputer(labelMessage string) string {
	validate := func(input string) error {
		return nil
	}
	prompt := promptui.Prompt{
		Label:    labelMessage,
		Validate: validate,
		Mask:     '*',
	}
	passwd, err := prompt.Run()
	if err != nil {
		log.Fatal("err@passwdInputer: Failed to run prompt")
	}

	return passwd
}
