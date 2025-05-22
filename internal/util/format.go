package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatRupiah(amount float64) string {
	formatted := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(formatted, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Tambahkan pemisah ribuan
	var result strings.Builder
	n := len(intPart)
	for i, digit := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	return result.String() + "," + decPart
}

func FormatTanggalIndo(dateStr string) (string, error) {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	bulan := [...]string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	tgl := parsed.Day()
	bln := bulan[int(parsed.Month())-1]
	thn := parsed.Year()

	return fmt.Sprintf("%02d %s %d", tgl, bln, thn), nil
}

func ParseRupiahToFloat64(rupiah string) (float64, error) {

	clean := strings.ReplaceAll(rupiah, ".", "")
	clean = strings.ReplaceAll(clean, ",", ".")
	return strconv.ParseFloat(clean, 64)
}
