package service

import (
	"errors"
	"sync"

	"photohttp/internal/images"
)

type PhotoService struct {
	img       *images.Client
	lastPhoto *images.Photo
	mu        sync.RWMutex
}

func NewPhotoService(img *images.Client) *PhotoService {
	return &PhotoService{img: img}
}

func (s *PhotoService) GetRandom() (*images.Photo, error) {
	photo, err := s.img.Random()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.lastPhoto = photo
	s.mu.Unlock()

	return photo, nil
}

func (s *PhotoService) Search(query string) (*images.Photo, error) {
	results, err := s.img.Search(query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no photos found")
	}

	s.mu.Lock()
	s.lastPhoto = &results[0]
	s.mu.Unlock()

	return &results[0], nil
}

func (s *PhotoService) GetLast() (*images.Photo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.lastPhoto == nil {
		return nil, errors.New("no last photo")
	}

	return s.lastPhoto, nil
}
