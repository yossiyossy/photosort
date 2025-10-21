package internal

import "strings"

// ▼ 判定系を「必ず内部で小文字に正規化」する実装に置き換え
func normExt(ext string) string {
	return strings.ToLower(ext)
}

func isImage(ext string) bool {
	ext = normExt(ext)
	return isJPEG(ext) || isRAW(ext)
}

func isJPEG(ext string) bool {
	switch normExt(ext) {
	case ".jpg", ".jpeg", ".jpe", ".jfif":
		return true
	default:
		return false
	}
}

func isRAW(ext string) bool {
	switch normExt(ext) {
	case ".cr2", ".cr3", ".nef", ".nrw", ".arw", ".srf", ".sr2",
		".orf", ".raf", ".rw2", ".rwl", ".dng", ".pef", ".srw":
		return true
	default:
		return false
	}
}
