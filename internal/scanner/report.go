package scanner

import (
	"fmt"
	"shadow-scanner/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDFReport - Natijalarni PDF faylga saqlaydi
func GeneratePDFReport(target string, results []models.ScanResult) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Sarlavha
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("ShadowScanner Hisoboti: %s", target))
	pdf.Ln(12)

	// Jadval sarlavhasi
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(30, 10, "Port", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 10, "Xizmat", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "Kechikish", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Holat", "1", 1, "C", true, 0, "")

	// Ma'lumotlar
	pdf.SetFont("Arial", "", 12)
	for _, res := range results {
		if res.State == "Open" {
			pdf.CellFormat(30, 10, fmt.Sprintf("%d", res.Port), "1", 0, "C", false, 0, "")
			pdf.CellFormat(60, 10, res.Service, "1", 0, "C", false, 0, "")
			pdf.CellFormat(50, 10, fmt.Sprintf("%v", res.Latency), "1", 0, "C", false, 0, "")
			pdf.SetTextColor(0, 150, 0) // Yashil rang
			pdf.CellFormat(40, 10, res.State, "1", 1, "C", false, 0, "")
			pdf.SetTextColor(0, 0, 0)
		}
	}

	return pdf.OutputFileAndClose("report.pdf")
}
