package media

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"users-api-v1/src/features"

	"github.com/gabriel-vasile/mimetype"
)

func PublicPrefix() string {
	publicPrefix := features.GetEnv("MEDIA_PUBLIC_PREFIX", "/media")
	return publicPrefix
}

var (
	ErrUnsupportedType = errors.New("Не поддерживаемый формат (Поддерживаеться: jpg, png, gif, webp, svg, heic, heif, avif, bmp, tiff)")
	ErrTooLarge        = errors.New("Файл слишком большой")
)

var allowed = []struct{ mime, ext string }{
	{"image/jpeg", ".jpg"},
	{"image/png", ".png"},
	{"image/gif", ".gif"},
	{"image/webp", ".webp"},
	{"image/svg+xml", ".svg"},
	{"image/heic", ".heic"},
	{"image/heif", ".heif"},
	{"image/avif", ".avif"},
	{"image/bmp", ".bmp"},
	{"image/tiff", ".tiff"},
}

func init() {
	_ = mime.AddExtensionType(".heic", "image/heic")
	_ = mime.AddExtensionType(".heif", "image/heif")
	_ = mime.AddExtensionType(".avif", "image/avif")
}

type Storage struct {
	dir     string
	maxSize int64
}

func NewStorage(dir string) (*Storage, error) {
	maxSize := features.GetEnvInt("MAX_AVATAR_SIZE_MB", 20)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Storage{dir: dir, maxSize: maxSize << 20}, nil
}

func (storage *Storage) MaxSize() int64 { return storage.maxSize }

func (storage *Storage) SaveImage(fileHeader *multipart.FileHeader) (string, error) {

	if fileHeader.Size > storage.maxSize {
		return "", fmt.Errorf("%w, поддерживается до %dмб", ErrTooLarge, storage.maxSize/(1<<20))
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mtype, err := mimetype.DetectReader(src)
	if err != nil {
		return "", err
	}
	ext, ok := extensionFor(mtype)

	if !ok {
		return "", ErrUnsupportedType
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	name, err := randomName(ext)
	if err != nil {
		return "", err
	}
	dstPath := filepath.Join(storage.dir, name)

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(dstPath)
		return "", err
	}
	if err := dst.Close(); err != nil {
		os.Remove(dstPath)
		return "", err
	}

	return PublicPrefix() + "/" + name, nil

}

func (storage *Storage) Delete(publicPath string) {
	if !IsLocal(publicPath) {
		return
	}
	name := filepath.Base(publicPath)
	err := os.Remove(filepath.Join(storage.dir, name))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("медиа: удалено %s: %v", name, err)
	}
}

func IsLocal(p string) bool {
	return strings.HasPrefix(p, PublicPrefix())
}

func extensionFor(m *mimetype.MIME) (string, bool) {
	for _, a := range allowed {
		if m.Is(a.mime) {
			return a.ext, true
		}
	}
	return "", false
}

func randomName(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}
