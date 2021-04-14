package main

import (
	"flag"
	"log"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

var lock = &sync.Mutex{}
var email = flag.String("e", "", "email")
var password = flag.String("p", "", "password")

func main() {
	flag.Parse()
	clockIn()
}

func clockIn() {
	if !isLoggedIn() {
		login()
	}

	lock.Lock()
	defer lock.Unlock()

	browser := newBrowser(true)
	defer browser.Close()

	page := browser.MustPage("https://www.futunn.com/account/home").MustWaitLoad()

	page.Race().Element("#signIn .none").MustHandle(func(el *rod.Element) {
		el.MustClick()
		log.Println("签到成功")
	}).ElementR("#signed", "已签到").MustHandle(func(el *rod.Element) {
		log.Println("已经签过到了")
	}).MustDo()
}

func isLoggedIn() bool {
	lock.Lock()
	defer lock.Unlock()

	browser := newBrowser(true)
	defer browser.Close()

	return browser.MustPage("https://www.futunn.com").MustWaitLoad().MustHasR("a", "退出|Sign Out")
}

func login() {
	lock.Lock()
	defer lock.Unlock()

	browser := newBrowser(false)
	defer browser.Close()

	page := browser.MustPage("https://passport.futunn.com/?target=https%3A%2F%2Fwww.futunn.com%2Faccount%2Fhome&lang=zh-hk#login")

	page.MustElement(`input[name=email]`).MustInput(*email)
	page.MustElement(`input[name=password]`).MustInput(*password).MustPress(input.Enter)

	page.Race().Element(".nn-header-account").MustHandle(func(e *rod.Element) {
		log.Println("登录成功，已经找到登录后页面元素内容 ", e.MustText())
	}).Element(".u-error-wrapper").MustHandle(func(e *rod.Element) {
		// 当用户名或密码错误时
		log.Println("失败")
		log.Panicln(e)
		log.Fatal(e.MustText())
	})
}

func newBrowser(headless bool) *rod.Browser {
	url := launcher.New().Headless(headless).UserDataDir("tmp/user").MustLaunch()
	return rod.New().ControlURL(url).MustConnect()
}
