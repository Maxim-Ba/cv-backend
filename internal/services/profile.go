package services

import (
	"fmt"
	"strings"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	"github.com/Maxim-Ba/cv-backend/internal/repository"
)

type ProfileReader interface {
	Get() (repository.Profile, error)
	GetAboutMeWithTechnologies() (repository.Profile, []dto.TechnologyWithTagsDTO, error)
	GetTechnologyIDs(profileID int64) ([]int64, error)
}

type ProfileAboutMeWriter interface {
	UpdateAboutMe(profileID int64, about, note, hobbies string) error
	SetTechnologies(profileID int64, technologyIDs []int64) error
}

type ProfileManager interface {
	ProfileReader
	ProfileAboutMeWriter
}

type ProfileService struct {
	repo ProfileManager
}

func NewProfileService(repo ProfileManager) *ProfileService {
	return &ProfileService{repo: repo}
}

// GetAboutMe возвращает данные секции «О себе» для API и админки
func (s *ProfileService) GetAboutMe() (dto.AboutMeDTO, error) {
	profile, technologies, err := s.repo.GetAboutMeWithTechnologies()
	if err != nil {
		return dto.AboutMeDTO{}, fmt.Errorf("get about me: %w", err)
	}

	return dto.AboutMeDTO{
		BioParagraphs: parseBioParagraphs(profile.About),
		Technologies:  technologies,
		Note:          profile.Note,
		Hobbies:       profile.Hobbies,
	}, nil
}

// UpdateAboutMeInput данные для обновления секции «О себе»
type UpdateAboutMeInput struct {
	About          string
	Note           string
	Hobbies        string
	TechnologyIDs  []int64
}

// UpdateAboutMe обновляет секцию «О себе» и связанные технологии
func (s *ProfileService) UpdateAboutMe(input UpdateAboutMeInput) error {
	profile, err := s.repo.Get()
	if err != nil {
		return fmt.Errorf("update about me: %w", err)
	}

	if err := s.repo.UpdateAboutMe(profile.ID, input.About, input.Note, input.Hobbies); err != nil {
		return fmt.Errorf("update about me: %w", err)
	}
	if err := s.repo.SetTechnologies(profile.ID, input.TechnologyIDs); err != nil {
		return fmt.Errorf("update about me technologies: %w", err)
	}
	return nil
}

func parseBioParagraphs(about *string) []string {
	if about == nil || strings.TrimSpace(*about) == "" {
		return []string{}
	}
	parts := strings.Split(*about, "\n\n")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// JoinBioParagraphs объединяет абзацы биографии для хранения в БД
func JoinBioParagraphs(paragraphs []string) string {
	filtered := make([]string, 0, len(paragraphs))
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return strings.Join(filtered, "\n\n")
}
