package config

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func login(config *Config) {
	// 创建上下文和取消函数
	ctx, _ := chromedp.NewExecAllocator(context.Background())

	ctx, _ = context.WithTimeout(ctx, 90*time.Second)
	ctx, _ = chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	if err := chromedp.Run(ctx, myTasks(config)); err != nil {
		log.Fatal(err)
		println(err.Error())
		return
	}

	chromedp.Cancel(ctx)

}

func myTasks(config *Config) chromedp.Tasks {
	return chromedp.Tasks{
		// 1. 打开金山文档的登陆界面
		chromedp.Navigate("https://account.dianping.com/pclogin"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 使用 chromedp 的 network.GetCookies 动作获取 Cookie
			total := 5
			for true {
				time.Sleep(1 * time.Second)
				total += 1
				cookies, err := network.GetCookies().Do(ctx)
				if err != nil {
					return err
				}

				// 根据需要处理 Cookie
				//var c []string
				for _, v := range cookies {
					if v.Name == "dper" {
						println("dper: ", v.Value)
						config.Dper = v.Value
						return nil
					}
					//aCookie := v.Name + " - " + v.Domain
					//c = append(c, aCookie)
				}

				//stringSlices := strings.Join(c[:], ",\n")
				//fmt.Printf("%v", stringSlices)
				if total == 50 {
					return nil
				}
			}
			return nil
		}),
	}
}
