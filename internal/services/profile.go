package services

import (
	"fmt"
	"strings"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	"github.com/Maxim-Ba/cv-backend/internal/repository"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

type ProfileReader interface {
	Get() (repository.Profile, error)
	GetAboutMeWithTechnologies(locale i18n.Locale) (repository.Profile, []dto.TechnologyWithTagsDTO, error)
	GetTechnologyIDs(profileID int64) ([]int64, error)
}

type ProfileAboutMeWriter interface {
	UpdateAboutMe(profileID int64, about, note, hobbies i18n.LocalizedText) error
	SetTechnologies(profileID int64, technologyIDs []int64) error
}

type ProfileHeroWriter interface {
	UpdateHero(profileID int64, greeting, fullName, title, pitch i18n.LocalizedText) error
}

type ProfileManager interface {
	ProfileReader
	ProfileAboutMeWriter
	ProfileHeroWriter
}

type ProfileService struct {
	repo ProfileManager
}

func NewProfileService(repo ProfileManager) *ProfileService {
	return &ProfileService{repo: repo}
}

// GetProfile возвращает профиль для админки
func (s *ProfileService) GetProfile() (repository.Profile, error) {
	return s.repo.Get()
}

// GetTechnologyIDsForAdmin возвращает ID технологий профиля
func (s *ProfileService) GetTechnologyIDsForAdmin() ([]int64, error) {
	profile, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	return s.repo.GetTechnologyIDs(profile.ID)
}

// GetHero возвращает локализованные данные hero-секции
func (s *ProfileService) GetHero(locale i18n.Locale) (dto.HeroDTO, error) {
	profile, err := s.repo.Get()
	if err != nil {
		return dto.HeroDTO{}, fmt.Errorf("get hero: %w", err)
	}
	return dto.HeroDTO{
		Greeting: profile.Greeting.Resolve(locale),
		FullName: profile.FullName.Resolve(locale),
		Title:    profile.Title.Resolve(locale),
		Pitch:    profile.Pitch.Resolve(locale),
	}, nil
}

// GetAboutMe возвращает данные секции «О себе» для API и админки
func (s *ProfileService) GetAboutMe(locale i18n.Locale) (dto.AboutMeDTO, error) {
	profile, technologies, err := s.repo.GetAboutMeWithTechnologies(locale)
	if err != nil {
		return dto.AboutMeDTO{}, fmt.Errorf("get about me: %w", err)
	}

	localizedTech := technologies

	return dto.AboutMeDTO{
		BioParagraphs: parseBioParagraphs(profile.About.Resolve(locale)),
		Technologies:  localizedTech,
		Note:          profile.Note.Resolve(locale),
		Hobbies:       profile.Hobbies.Resolve(locale),
	}, nil
}

// GetAboutMeAdmin возвращает данные «О себе» для админки с RU/EN значениями
func (s *ProfileService) GetAboutMeAdmin() (repository.Profile, []dto.TechnologyWithTagsDTO, error) {
	return s.repo.GetAboutMeWithTechnologies(i18n.LocaleRU)
}

// UpdateAboutMeInput данные для обновления секции «О себе»
type UpdateAboutMeInput struct {
	About         i18n.LocalizedText
	Note          i18n.LocalizedText
	Hobbies       i18n.LocalizedText
	TechnologyIDs []int64
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

// UpdateHeroInput данные для обновления hero-секции
type UpdateHeroInput struct {
	Greeting i18n.LocalizedText
	FullName i18n.LocalizedText
	Title    i18n.LocalizedText
	Pitch    i18n.LocalizedText
}

// UpdateHero обновляет hero-секцию профиля
func (s *ProfileService) UpdateHero(input UpdateHeroInput) error {
	profile, err := s.repo.Get()
	if err != nil {
		return fmt.Errorf("update hero: %w", err)
	}
	if err := s.repo.UpdateHero(profile.ID, input.Greeting, input.FullName, input.Title, input.Pitch); err != nil {
		return fmt.Errorf("update hero: %w", err)
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
