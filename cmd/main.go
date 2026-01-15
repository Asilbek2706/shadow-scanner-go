package main

import (
	"fmt"
	"net/http"
	"os" // Muhit o'zgaruvchilari bilan ishlash uchun
	"shadow-scanner/internal/models"
	"shadow-scanner/internal/scanner"
	"shadow-scanner/internal/ws"
	"sync"
	"time"
)

// Global o'zgaruvchilar hisobot uchun natijalarni vaqtinchalik saqlaydi
var (
	lastResults []models.ScanResult
	lastTarget  string
	resultsMu   sync.Mutex
)

func main() {
	hub := ws.NewHub()

	// Static fayllar uchun xizmat (index.html va h.k.)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// WebSocket handler
	http.HandleFunc("/ws", hub.HandleWS)

	// PDF yuklab olish uchun yo'l
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
		w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, "report.pdf")
	})

	// Skanerlashni boshlash API
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("target")
		if target == "" {
			target = "scanme.nmap.org"
		}

		resultsMu.Lock()
		lastTarget = target
		lastResults = []models.ScanResult{}
		resultsMu.Unlock()

		go func() {
			ports := make(chan int, 100)
			results := make(chan models.ScanResult)
			pool := scanner.Pool{WorkerCount: 100, Timeout: 500 * time.Millisecond}

			go func() {
				pool.Start(target, ports, results)
				close(results)
			}()

			go func() {
				for i := 1; i <= 1024; i++ {
					ports <- i
				}
				close(ports)
			}()

			for res := range results {
				if res.State == "Open" {
					resultsMu.Lock()
					lastResults = append(lastResults, res)
					resultsMu.Unlock()
					hub.Broadcast(res)
				}
			}

			fmt.Printf("Skanerlash tugadi: %s. Hisobot tayyorlanmoqda...\n", target)
			scanner.GeneratePDFReport(lastTarget, lastResults)
			hub.Broadcast(map[string]string{"status": "completed"})
		}()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Scan started")
	})

	// Render portini olish
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Lokal muhit uchun default port
	}

	fmt.Printf("🚀 ShadowScanner Dashboard ishga tushdi: http://localhost:%s\n", port)
	
	// Serverni ishga tushirish
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Xatolik yuz berdi: %v\n", err)
	}
}
