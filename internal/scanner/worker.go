package scanner

import (
	"shadow-scanner/internal/models"
	"sync"
	"time"
)

// Pool skanerlash ishchilarining sozlamalarini saqlaydi
type Pool struct {
	WorkerCount int
	Timeout     time.Duration
}

// Start ishchilarni (workers) ishga tushiradi
func (p *Pool) Start(ip string, ports <-chan int, results chan<- models.ScanResult) {
	var wg sync.WaitGroup

	for i := 0; i < p.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Kanaldan yangi port raqami kelishini kutadi
			for port := range ports {
				res := ScanPort(ip, port, p.Timeout)
				results <- res
			}
		}()
	}

	wg.Wait()
}
