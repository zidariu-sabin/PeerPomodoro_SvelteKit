package domain

import (
	"sync"
	"time"
)

type Timer struct {
	WorkTime         int64 `json:"work_time"`
	BreakTime        int64 `json:"break_time"`
	TotalRounds      int64 `json:"total_rounds"`
	CurrentRound     int64 `json:"current_round"`
	SecondsRemaining int64 `json:"seconds_remaining"`
	IsWorkPeriod     bool  `json:"is_work_period"`
	IsRunning        bool  `json:"is_running"`
	IsCompleted      bool  `json:"is_completed"`
	IsTimerCreated   bool  `json:"is_timer_created"`

	ticker *time.Ticker
	done   chan struct{}
	mu     sync.Mutex
}

func NewTimer(workTime int64, breakTime int64, totalRounds int64) *Timer {
	return &Timer{
		WorkTime:         workTime,
		BreakTime:        breakTime,
		TotalRounds:      totalRounds,
		CurrentRound:     0,
		SecondsRemaining: workTime * 60,
		IsWorkPeriod:     true,
		IsRunning:        false,
		IsCompleted:      false,
		IsTimerCreated:   true,
		ticker:           nil,
	}
}

func (t *Timer) InitializeTimer() {
	t.IsTimerCreated = true
	t.CurrentRound = 1
	t.IsCompleted = false
	t.SecondsRemaining = t.WorkTime * 60
	t.IsWorkPeriod = true
	t.IsRunning = false
}

func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.IsRunning {
		return
	}

	t.IsRunning = false
	if t.ticker != nil {
		t.ticker.Stop()
	}
	close(t.done)
}

func (t *Timer) tick() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.SecondsRemaining > 0 {
		t.SecondsRemaining--
	} else {
		t.handlePeriodComplete()
	}
}

func (t *Timer) Start() {
	t.mu.Lock()
	if t.IsRunning {
		t.mu.Unlock()
	}
	t.IsRunning = true
	t.ticker = time.NewTicker(time.Second)
	t.done = make(chan struct{})
	t.mu.Unlock()
	go func() {
		for t.CurrentRound <= t.TotalRounds {
			select {
			case <-t.done:
				return
			case <-t.ticker.C:
				t.tick()
			}
			if t.SecondsRemaining > 0 {
				t.SecondsRemaining--
				//relayed time update
				// minutes := t.SecondsRemaining / 60
				// seconds := t.SecondsRemaining % 60
			} else {
				t.handlePeriodComplete()
			}
		}
	}()
}
func (t *Timer) handlePeriodComplete() {
	if t.IsWorkPeriod {
		t.IsWorkPeriod = false
		t.SecondsRemaining = t.BreakTime * 60
	} else {
		t.IsWorkPeriod = true
		if t.CurrentRound >= t.TotalRounds {
			// pause()
			t.IsCompleted = true
			t.IsRunning = false
			if t.ticker != nil {
				t.ticker.Stop()
			}
			close(t.done)
			return
		}
		t.CurrentRound++
		t.SecondsRemaining = t.WorkTime * 60
	}
}
