package image

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"syscall"
)

type ImageController struct {
	imageStoragePath string
}

func NewImageController(imageStoragePath string) (*ImageController, error) {
	if err := os.MkdirAll(imageStoragePath, 0755); err != nil {
		return nil, err

	}
	return &ImageController{imageStoragePath: imageStoragePath}, nil
}

func (c *ImageController) isNeedRefreshImage(url string, targetFile string) (bool, error) {
	// todo смотреть на заголовки
	file, err := os.OpenFile(targetFile, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return false, err
	}
	defer file.Close()

	return false, nil
}

func (c *ImageController) DownloadImage(url string, name string, force bool) error {
	targetFile := path.Join(c.imageStoragePath, name)

	isNeedRefreshImage, err := c.isNeedRefreshImage(url, targetFile)
	if !force && !isNeedRefreshImage && err == nil {
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cant make http get: %w", err)
	}
	defer response.Body.Close()

	tmpFile, err := ioutil.TempFile("", "image_"+name)
	if err != nil {
		return fmt.Errorf("cant create tmp tmpFile for image: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, response.Body); err != nil {
		return fmt.Errorf("cant download content: %w", err)
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("cant seek tmp file: %w", err)
	}

	file, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("cant create target file (%s): %w", targetFile, err)
	}
	defer file.Close()

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("cant flock target file (%s): %w", targetFile, err)

	}
	if _, err := io.Copy(file, tmpFile); err != nil {
		return fmt.Errorf("cant download content: %w", err)
	}

	return nil
}

func (c *ImageController) UnpackImage(name string, targetPath string) error {

	if _, err := os.Open(targetPath); err != nil {
		return fmt.Errorf("cant open targetPath dir (%s): %w", targetPath, err)
	}

	err := os.RemoveAll(targetPath)
	if err != nil {
		return fmt.Errorf("cant clean targetPath (%s): %w", targetPath, err)
	}

	targz := archiver.NewTarGz()
	source := path.Join(c.imageStoragePath, name)

	if err := targz.Unarchive(source, targetPath); err != nil {
		return fmt.Errorf("cant unarchive %s: %w", source, err)
	}
	defer targz.Close()

	return nil
}
