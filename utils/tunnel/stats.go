package tunnel

import (
	"sync"
	"time"
)

type Stats struct {
	*sync.Mutex
	State               State
	currentTxCount      int64
	currentRxCount      int64
	TxCount             int64
	RxCount             int64
	TxSpeed             int64
	RxSpeed             int64
	ConnQty             int64
	UpdateStatsCallback func(stats *Stats)
}

func NewStats(state State, updateStatsCallback func(stats *Stats)) *Stats {
	stats := &Stats{
		Mutex:               &sync.Mutex{},
		State:               state,
		UpdateStatsCallback: updateStatsCallback,
	}

	if updateStatsCallback != nil {
		stats.InitSpeedCalculation()
	}

	return stats
}

func (s *Stats) InitSpeedCalculation() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.calculateSpeed()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Stats) AddRxCount(count int64) {
	s.Lock()

	s.RxCount += count
	s.currentRxCount += count

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) AddTxCount(count int64) {
	s.Lock()

	s.TxCount += count
	s.currentTxCount += count

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) AddConnQty() {
	s.Lock()

	s.ConnQty++

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) SubConnQty() {
	s.Lock()

	if s.ConnQty > 0 {
		s.ConnQty--
	}

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) SetConnQty(qty int64) {
	s.Lock()

	s.ConnQty = qty

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) UpdateState(state State) {
	s.Lock()

	s.State = state

	s.Unlock()
	s.triggerUpdateStatsCallback()
}

func (s *Stats) triggerUpdateStatsCallback() {
	if s.UpdateStatsCallback != nil {
		s.UpdateStatsCallback(s)
	}
}

func (s *Stats) calculateSpeed() {
	s.Lock()

	s.TxSpeed = s.currentTxCount
	s.currentTxCount = 0

	s.RxSpeed = s.currentRxCount
	s.currentRxCount = 0

	s.Unlock()
	s.triggerUpdateStatsCallback()
}
