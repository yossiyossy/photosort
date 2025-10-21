// internal/exifutil.go
package internal

import (
	"fmt"
	"os"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// GetDateFromExif は画像から撮影日時を推定して返します。
// 優先順位: DateTimeOriginal -> CreateDate -> ファイルのModTime
func GetDateFromExif(path string) (time.Time, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("read file: %w", err)
	}

	// 1) 画像バイト列から EXIF ブロックを抽出
	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		// EXIF自体が無い場合はFSの時刻にフォールバック
		fi, statErr := os.Stat(path)
		if statErr != nil {
			return time.Time{}, fmt.Errorf("no exif and stat failed: %v / %v", err, statErr)
		}
		return fi.ModTime(), nil
	}

	// 2) IFD/タグの定義をロード（標準タグセット）
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return time.Time{}, fmt.Errorf("ifd mapping: %w", err)
	}
	ti := exif.NewTagIndex()

	// 3) フラットなタグ一覧を取得（この方法が一番シンプル）
	tags, _, err := exif.GetFlatExifData(rawExif, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("get flat exif: %w", err)
	}
	_ = im
	_ = ti // （今は未使用だが Collect 等に切り替える時に使う）

	// 4) 欲しい日時タグを探す
	var val string
	for _, name := range []string{"DateTimeOriginal", "CreateDate"} {
		for _, t := range tags {
			if t.TagName == name && t.FormattedFirst != "" {
				val = t.FormattedFirst
				break
			}
		}
		if val != "" {
			break
		}
	}

	// 5) 見つかった文字列を EXIF 形式でパース
	if val != "" {
		// EXIFの典型フォーマット: "2006:01:02 15:04:05"
		parsed, perr := time.Parse("2006:01:02 15:04:05", val)
		if perr == nil {
			return parsed, nil
		}
		// 稀に "2006-01-02 15:04:05" のこともあるので緩和
		if parsed2, perr2 := time.Parse("2006-01-02 15:04:05", val); perr2 == nil {
			return parsed2, nil
		}
	}

	// 6) ここまでで取れなければ FS の時刻を返す
	fi, statErr := os.Stat(path)
	if statErr != nil {
		return time.Time{}, fmt.Errorf("fallback stat failed: %w", statErr)
	}
	return fi.ModTime(), nil
}
