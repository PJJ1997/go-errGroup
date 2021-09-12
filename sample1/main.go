package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
*/

func main() {
	// ============================ 示例一 ============================
	// 验证使用errGroup创建的goroutine不会因为其中一个goroutine出错而取消执行
	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < 3; i++ {
		if i == 1 {
			g.Go(func() error {
				return errors.New("访问redis失败\n")
			})
		} else {
			g.Go(func() error {
				time.Sleep(5 * time.Second)
				fmt.Println("i am fine")
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("appear error and err is %s", err.Error())
	}

	time.Sleep(2 * time.Second)
}
