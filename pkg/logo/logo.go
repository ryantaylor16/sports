package logo

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	log "github.com/sirupsen/logrus"

	rgb "github.com/robbydyer/sports/pkg/rgbmatrix-rpi"
	"github.com/robbydyer/sports/pkg/rgbrender"
)

// Logo is used to manage logo rendering
type Logo struct {
	key             string
	sourceLogo      image.Image
	bounds          image.Rectangle
	targetDirectory string
	config          *Config
	thumbnail       image.Image
	log             *log.Logger
}

// Config ...
type Config struct {
	Abbrev string `json:"abbrev"`
	Pt     *Pt    `json:"pt"`
}

// Pt defines the x, y shift and zoom values for a logo
type Pt struct {
	X    int     `json:"xShift"`
	Y    int     `json:"yShift"`
	Zoom float64 `json:"zoom"`
}

// New ...
func New(key string, sourceLogo image.Image, targetDirectory string, matrixBounds image.Rectangle, conf *Config) *Logo {
	return &Logo{
		key:             key,
		targetDirectory: targetDirectory,
		sourceLogo:      sourceLogo,
		config:          conf,
		bounds:          matrixBounds,
	}
}

func (l *Logo) SetLogger(logger *log.Logger) {
	l.log = logger
}

func (l *Logo) ensureLogger() {
	if l.log == nil {
		l.log = log.New()
	}
}

// ThumbnailFilename returns the filname for the resized thumbnail to use
func (l *Logo) ThumbnailFilename(size image.Rectangle) string {
	return fmt.Sprintf("%s/%s_%dx%d_%f.png", l.targetDirectory, l.key, size.Dx(), size.Dy(), l.config.Pt.Zoom)
}

// GetThumbnail returns the resized image
func (l *Logo) GetThumbnail(size image.Rectangle) (image.Image, error) {
	if l.thumbnail != nil {
		return l.thumbnail, nil
	}

	thumbFile := l.ThumbnailFilename(size)

	if _, err := os.Stat(thumbFile); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Stat(l.targetDirectory); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(l.targetDirectory, 0755); err != nil {
						return nil, fmt.Errorf("failed to create logo cache dir: %w", err)
					}
				}
			}
			// Create the thumbnail
			l.thumbnail = rgbrender.ResizeImage(l.sourceLogo, size, l.config.Pt.Zoom)

			go func() {
				l.ensureLogger()
				l.log.Infof("Saving thumbnail logo %s\n", thumbFile)
				if err := rgbrender.SavePng(l.thumbnail, thumbFile); err != nil {
					l.log.Error("Failed to save logo PNG", err)
				}
			}()

			return l.thumbnail, nil
		}

		return nil, err
	}

	t, err := os.Open(thumbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open logo %s: %w", thumbFile, err)
	}
	defer t.Close()

	l.thumbnail, err = png.Decode(t)
	if err != nil {
		return nil, fmt.Errorf("failed to decode logo %s: %w", thumbFile, err)
	}

	return l.thumbnail, nil
}

// RenderLeftAligned renders the logo on the left side of the matrix
func (l *Logo) RenderLeftAligned(canvas *rgb.Canvas, width int) (image.Image, error) {
	thumb, err := l.GetThumbnail(l.bounds)
	if err != nil {
		return nil, err
	}

	startX := width - thumb.Bounds().Dx() + l.config.Pt.X
	startY := 0 + l.config.Pt.Y

	bounds := image.Rect(startX, startY, canvas.Bounds().Dx()-1, canvas.Bounds().Dy()-1)

	i := image.NewRGBA(bounds)
	draw.Draw(i, bounds, thumb, image.Point{}, draw.Over)

	return i, nil
}

// RenderRightAligned renders the logo on the right side of the matrix
func (l *Logo) RenderRightAligned(canvas *rgb.Canvas, width int) (image.Image, error) {
	thumb, err := l.GetThumbnail(l.bounds)
	if err != nil {
		return nil, err
	}

	startX := width + l.config.Pt.X
	startY := 0 + l.config.Pt.Y

	bounds := image.Rect(startX, startY, canvas.Bounds().Dx()-1, canvas.Bounds().Dy()-1)

	i := image.NewRGBA(bounds)
	draw.Draw(i, bounds, thumb, image.Point{}, draw.Over)

	return i, nil
}
