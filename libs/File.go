package libs

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Width     int
	Height    int
	Crop      bool
	Quality   int
	Rotate    int
	Interlace bool
	image     image.Image
}

//SaveImage
func (file *File) SaveImage(c *gin.Context, image *multipart.FileHeader, dirname string) (directory string, err error) {
	filename := filepath.Base(image.Filename)
	directory = fmt.Sprintf("%s/%s", dirname, filename)
	if err = c.SaveUploadedFile(image, directory); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	return
}
func (file *File) New(width int, height int, quality int) {
	file.Width = width
	file.Height = height
	file.Quality = quality
	file.Crop = true
	file.Interlace = true
}
func (file *File) Read(filename string) {
	imagePath, err := imaging.Open(filename)
	file.image = imagePath
	if err != nil {
		log.Fatalf("failed to open imagePath: %v", err)
	}
}
func (file *File) Resize() {
	file.image = imaging.Resize(file.image, file.Width, file.Height, imaging.Lanczos)
}
func (file *File) SaveResize(dir, filename string) (name string, err error) {
	var directory = fmt.Sprintf("%s/%s", dir, filename)
	dst := imaging.Resize(file.image, file.Width, file.Height, imaging.Lanczos)
	err = imaging.Save(dst, directory, imaging.JPEGQuality(file.Quality))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	name = filename
	return
}

func (file *File) GenerateName(s string) string {
	//n := 20
	//b := make([]byte, n)
	//if _, err := rand.Read(b); err != nil {
	//	panic(err)
	//}
	//result := fmt.Sprintf("%X", b)
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)))
	return fmt.Sprintf("%s.%s", formatted, s)
}

func (file *File) DeleteFile(path string) {
	var err = os.Remove(path)
	if err != nil {
		log.Fatalf("failed to delete image: %v", err)
	}
	fmt.Println("File Deleted")
}
