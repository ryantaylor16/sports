package espnboard

import (
	"context"
	"path/filepath"

	"go.uber.org/zap"
)

type nfl struct{}

func (n *nfl) League() string {
	return "NFL"
}

func (n *nfl) APIPath() string {
	return "football/nfl"
}

func (n *nfl) TeamEndpoints() []string {
	return []string{filepath.Join(n.APIPath(), "groups")}
}

func (n *nfl) HTTPPathPrefix() string {
	return "nfl"
}

// NewNFL ...
func NewNFL(ctx context.Context, logger *zap.Logger) (*ESPNBoard, error) {
	return New(ctx, &nfl{}, logger)
}

type ncaam struct{}

func (n *ncaam) League() string {
	return "NCAA Basketball"
}

func (n *ncaam) APIPath() string {
	return "basketball/mens-college-basketball"
}

func (n *ncaam) TeamEndpoints() []string {
	return []string{
		filepath.Join(n.APIPath(), "groups"),
		filepath.Join(n.APIPath(), "teams"),
	}
}

func (n *ncaam) HTTPPathPrefix() string {
	return "ncaam"
}

// NewNCAAMensBasketball ...
func NewNCAAMensBasketball(ctx context.Context, logger *zap.Logger) (*ESPNBoard, error) {
	return New(ctx, &ncaam{}, logger)
}

type nba struct{}

func (n *nba) League() string {
	return "NBA"
}

func (n *nba) APIPath() string {
	return "basketball/nba"
}

func (n *nba) TeamEndpoints() []string {
	return []string{filepath.Join(n.APIPath(), "groups")}
}

func (n *nba) HTTPPathPrefix() string {
	return "nba"
}

// NewNBA ...
func NewNBA(ctx context.Context, logger *zap.Logger) (*ESPNBoard, error) {
	return New(ctx, &nba{}, logger)
}

type mls struct{}

func (n *mls) League() string {
	return "MLS"
}

func (n *mls) APIPath() string {
	return "soccer/usa.1"
}

func (n *mls) TeamEndpoints() []string {
	return []string{filepath.Join(n.APIPath(), "teams")}
}

func (n *mls) HTTPPathPrefix() string {
	return "mls"
}

// NewMLS ...
func NewMLS(ctx context.Context, logger *zap.Logger) (*ESPNBoard, error) {
	return New(ctx, &mls{}, logger)
}
