/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 23:32
 */

package task

import "time"

type TimerFunc func(interface{}) bool

/**
 * 定时调用
 * @delay 首次延时
 * @tick  间隔
 * @fun   定时执行function
 * @param fun参数
 */

// Timer 定时调用
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				if fun(param) == false {
					return
				}
				t.Reset(tick)
			}
		}
	}()
}
