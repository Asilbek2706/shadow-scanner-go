package scanner

import (
	"fmt"
	"net"
	"shadow-scanner/internal/models"
	"time"
)

// ScanPort berilgan IP va portga TCP ulanish hosil qilib ko'radi
func ScanPort(ip string, port int, timeout time.Duration) models.ScanResult {
	address := fmt.Sprintf("%s:%d", ip, port)
	result := models.ScanResult{Port: port}

	start := time.Now()

	// Tarmoq orqali ulanishga urinish
	conn, err := net.DialTimeout("tcp", address, timeout)
	result.Latency = time.Since(start)

	if err != nil {
		result.State = "Closed"
		return result
	}

	conn.Close()
	result.State = "Open"
	result.SetService()
	return result
}
