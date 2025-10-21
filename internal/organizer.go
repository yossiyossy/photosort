package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// OrganizeInPlace は、exeDir 直下のファイルを走査し、JPG/RAW に分類しつつ
// EXIF日時(YYYYMMDD)のディレクトリを切って移動します。
func OrganizeInPlace(exeDir string) error {
	entries, err := os.ReadDir(exeDir)
	if err != nil {
		return fmt.Errorf("readdir: %w", err)
	}

	for _, e := range entries {
		// ディレクトリには触らない
		if e.IsDir() {
			continue
		}

		src := filepath.Join(exeDir, e.Name())

		ext := filepath.Ext(e.Name())
		if !isImage(ext) {
			continue
		}

		var top string
		if isJPEG(ext) {
			top = "JPG"
		} else if isRAW(ext) {
			top = "RAW"
		} else {
			continue
		}

		// EXIF日時（なければファイルのModTime）を取得
		dt, derr := GetDateFromExif(src)
		if derr != nil || dt.IsZero() {
			fi, statErr := os.Stat(src)
			if statErr != nil {
				fmt.Printf("skip (stat failed): %s (%v)\n", src, statErr)
				continue
			}
			dt = fi.ModTime()
		}
		dateFolder := dt.In(time.Local).Format("20060102") // YYYYMMDD

		dstDir := filepath.Join(exeDir, top, dateFolder)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		dst := filepath.Join(dstDir, e.Name())

		// 既に同名があるならスキップ
		if _, err := os.Stat(dst); err == nil {
			fmt.Printf("skip (already exists): %s\n", dst)
			continue
		}

		// パーティションをまたぐことがないのでRenameでOK
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("rename: %w", err)
		}

		fmt.Printf("%s → %s\n", e.Name(), filepath.Join(top, dateFolder))
	}

	return nil
}
