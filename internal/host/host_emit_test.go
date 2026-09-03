package host

import (
	"strings"
	"sync"
	"testing"
)

func TestVietnameseHostOutputBlocksUnexpectedHanText(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")
	h := &Host{events: make(chan Event, 1), streamCh: make(chan string, 1), done: make(chan struct{}, 1)}

	h.emitEvent(Event{Summary: "恢复创作", Detail: "查询上下文"})
	h.emitDelta("正文夹杂中文")

	ev := <-h.events
	delta := <-h.streamCh
	if strings.ContainsAny(ev.Summary+ev.Detail+delta, "恢复创作查询上下文正文夹杂中文") {
		t.Fatalf("Han text reached Vietnamese output: event=%+v delta=%q", ev, delta)
	}
}

// 关闭后的 emit 应被明确拒绝，不能依赖 recover 吞掉竞态。
func TestEmitAfterCloseDoesNotPanic(t *testing.T) {
	h := &Host{
		events:   make(chan Event, 1),
		streamCh: make(chan string, 1),
		done:     make(chan struct{}, 1),
	}
	h.closeOutputChannels()

	h.emitEvent(Event{Summary: "after close"})
	h.emitDelta("after close")
}

func TestConcurrentEmitAndCloseDoesNotRaceChannelLifecycle(t *testing.T) {
	h := &Host{
		events:   make(chan Event, 1),
		streamCh: make(chan string, 1),
		done:     make(chan struct{}, 1),
	}

	var emitters sync.WaitGroup
	for range 8 {
		emitters.Add(1)
		go func() {
			defer emitters.Done()
			for range 100 {
				h.emitEvent(Event{})
				h.emitDelta("delta")
			}
		}()
	}
	h.closeOutputChannels()
	emitters.Wait()

	if !h.outputClosed {
		t.Fatal("closeOutputChannels 应原子标记输出已关闭")
	}
}
