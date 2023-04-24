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
	//ctx, _ := chromedp.NewContext(context.Background()) /*
	// level0-debug1: Instance created with logs
	//ctx, _ := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf)) /*
	// level0-debug2: Create instances in no-headless mode
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("enable-automation", false),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//*/

	// level1: Access the page
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.Click(`//*[@id="gb"]/div/div[2]/a`, chromedp.NodeVisible),
	); err != nil {
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
		log.Fatal("err3-1@login: Failed login")
	}
	time.Sleep(10 * time.Second)

	// level3-2: Page Transition Confirmation
	var url2CheckTransition string
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`//*[@id="headingText"]/span`),
		chromedp.Location(&url2CheckTransition),
	); err != nil {
		log.Fatal("err3-2@login: Failed in page transition confirmation process")
	}
	if strings.Contains(url2CheckTransition, "dsh") {
		fmt.Println(url2CheckTransition)
		log.Fatal("err3-3@login: Failed to load on email address input page")
	}

	// level4: Processing password
	if err := chromedp.Run(ctx,
		input.InsertText(password),
		chromedp.Click(`//*[@id="passwordNext"]/div/button/div[3]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`//*[@id="headingText"]/span`),
	); err != nil {
		log.Fatal("err4-1@login: Failed to operate the login button")
	}
	time.Sleep(5 * time.Second)

	// level4-2: Page Transition Confirmation
	if err := chromedp.Run(ctx,
		chromedp.Location(&url2CheckTransition),
	); err != nil {
		log.Fatal("err4-2@login: Failed in page transition confirmation process")
	}
	if strings.Contains(url2CheckTransition, "pwd") {
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
