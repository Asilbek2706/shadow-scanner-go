package models

import "time"

// ScanResult har bir port tekshiruvi natijasini saqlaydi
type ScanResult struct {
	Port    int           `json:"port"`
	State   string        `json:"state"`   // "Open" yoki "Closed"
	Service string        `json:"service"` // Masalan: "HTTP", "SSH"
	Latency time.Duration `json:"latency"`
}

// SetService port raqamiga qarab xizmat turini taxmin qiladi
func (r *ScanResult) SetService() {
	commonPorts := map[int]string{
		21:   "FTP",
		22:   "SSH",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		443:  "HTTPS",
		3306: "MySQL",
		5432: "PostgreSQL",
		6379: "Redis",
		8080: "HTTP-Proxy",
	}

	if service, ok := commonPorts[r.Port]; ok {
		r.Service = service
	} else {
		r.Service = "Unknown"
	}
}
