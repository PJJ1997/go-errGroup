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
	// ============================ 示例二 ============================
	// 验证使用WithContext创建的errGroup其中一个goroutine出错，子goroutine会被cancel取消执行
	g, ctx := errgroup.WithContext(context.Background())

	// 单独开一个协程去做其他的事情，不参与waitGroup
	go WriteChangeLog(ctx)

	for i := 0; i < 3; i++ {
		g.Go(func() error {
			return errors.New("访问redis失败\n")
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("appear error and err is %s", err.Error())
	}
	time.Sleep(1 * time.Second)
}

func WriteChangeLog(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case <-time.After(time.Millisecond * 50):
		fmt.Println("write changelog")
	}
	return nil
}
