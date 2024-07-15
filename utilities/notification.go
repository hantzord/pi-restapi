package utilities

import (
	"capstone/constants"
	"fmt"
)

func AddContentComplaintUserNotification(name, message string) string {
	content := fmt.Sprintf("User %s has added a new complaint with message: %s", name, message)
	return content
}

func AddContentConsultationUserNotification(name, status string) string {
	switch status {
	case constants.REJECTED:
		return fmt.Sprintf("Kami mohon maaf, tetapi kami tidak dapat menjadwalkan sesi konsultasi Anda bersama dokter %s saat ini. Harap tunggu konfirmasi lebih lanjut", name)
	case constants.INCOMING:
		return fmt.Sprintf("Waktunya untuk konsultasi! Jadwalkan sesi dengan dokter %s Anda sekarang untuk mendapatkan dukungan yang Anda butuhkan.", name)
	case constants.PENDING:
		return fmt.Sprintf("Kami telah menerima permintaan konsultasi Anda. Harap tunggu konfirmasi lebih lanjut dari dokter %s", name)
	case constants.DONE:
		return fmt.Sprintf("Sesi konsultasi Anda bersama dokter %s telah selesai. Terima kasih telah menggunakan layanan kami", name)
	default:
		return "konsultasi Anda telah selesai. Terima kasih telah menggunakan layanan kami"
	}
}
