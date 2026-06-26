package repository

import (
	"database/sql"
	"fmt"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
)

// Profile хранит контактные данные владельца CV
type Profile struct {
	ID       int64
	FullName string
	Title    string
	About    *string
	Note     *string
	Hobbies  *string
	Email    *string
	Telegram *string
	GitHub   *string
	Phone    *string
}

// ProfileRepo репозиторий для чтения и обновления профиля
type ProfileRepo struct {
	db *sql.DB
}

// NewProfileRepo создает новый экземпляр ProfileRepo
func NewProfileRepo(db *sql.DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

// Get возвращает единственную запись профиля
func (r *ProfileRepo) Get() (Profile, error) {
	const query = `
		SELECT id, full_name, title, about, note, hobbies, email, telegram, github, phone
		FROM profile
		LIMIT 1
	`
	var p Profile
	err := r.db.QueryRow(query).Scan(
		&p.ID,
		&p.FullName,
		&p.Title,
		&p.About,
		&p.Note,
		&p.Hobbies,
		&p.Email,
		&p.Telegram,
		&p.GitHub,
		&p.Phone,
	)
	if err == sql.ErrNoRows {
		return Profile{}, fmt.Errorf("profile not found")
	}
	if err != nil {
		return Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	return p, nil
}

// GetAboutMeWithTechnologies возвращает данные секции «О себе» с технологиями-бейджами
func (r *ProfileRepo) GetAboutMeWithTechnologies() (Profile, []dto.TechnologyWithTagsDTO, error) {
	profile, err := r.Get()
	if err != nil {
		return Profile{}, nil, err
	}

	const query = `
		SELECT t.id, t.title, t.description, t.logo_url,
		       tg.id, tg.name, tg.hex_color
		FROM profile_technology pt
		JOIN technology t ON t.id = pt.technology_id
		LEFT JOIN technologies_tag tt ON tt.technology_id = t.id
		LEFT JOIN tag tg ON tg.id = tt.tag_id
		WHERE pt.profile_id = $1
		ORDER BY t.title, tg.id
	`

	rows, err := r.db.Query(query, profile.ID)
	if err != nil {
		return Profile{}, nil, fmt.Errorf("failed to query profile technologies: %w", err)
	}
	defer rows.Close()

	techMap := make(map[int64]*dto.TechnologyWithTagsDTO)
	techOrder := []int64{}

	for rows.Next() {
		var (
			techID      int64
			techTitle   string
			techDesc    sql.NullString
			techLogo    sql.NullString
			tagID       sql.NullInt64
			tagName     sql.NullString
			tagHexColor sql.NullString
		)
		if err := rows.Scan(
			&techID, &techTitle, &techDesc, &techLogo,
			&tagID, &tagName, &tagHexColor,
		); err != nil {
			return Profile{}, nil, fmt.Errorf("failed to scan profile technology row: %w", err)
		}

		if _, exists := techMap[techID]; !exists {
			tech := &dto.TechnologyWithTagsDTO{
				TechnologyDTO: dto.TechnologyDTO{
					ID:    techID,
					Title: techTitle,
				},
				Tags: []dto.TagDTO{},
			}
			if techDesc.Valid {
				tech.Description = &techDesc.String
			}
			if techLogo.Valid {
				tech.LogoUrl = &techLogo.String
			}
			techMap[techID] = tech
			techOrder = append(techOrder, techID)
		}

		if tagID.Valid {
			techMap[techID].Tags = append(techMap[techID].Tags, dto.TagDTO{
				ID:       tagID.Int64,
				Name:     tagName.String,
				HexColor: tagHexColor.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return Profile{}, nil, fmt.Errorf("profile technology rows error: %w", err)
	}

	technologies := make([]dto.TechnologyWithTagsDTO, 0, len(techOrder))
	for _, id := range techOrder {
		technologies = append(technologies, *techMap[id])
	}

	return profile, technologies, nil
}

// UpdateAboutMe обновляет текстовые поля секции «О себе»
func (r *ProfileRepo) UpdateAboutMe(profileID int64, about, note, hobbies string) error {
	const query = `
		UPDATE profile
		SET about = NULLIF($2, ''), note = NULLIF($3, ''), hobbies = NULLIF($4, '')
		WHERE id = $1
	`
	result, err := r.db.Exec(query, profileID, about, note, hobbies)
	if err != nil {
		return fmt.Errorf("failed to update profile about me: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("profile not found")
	}
	return nil
}

// SetTechnologies заменяет набор технологий-бейджей у профиля
func (r *ProfileRepo) SetTechnologies(profileID int64, technologyIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec("DELETE FROM profile_technology WHERE profile_id = $1", profileID); err != nil {
		return fmt.Errorf("failed to remove existing profile technologies: %w", err)
	}

	for _, technologyID := range technologyIDs {
		if _, err = tx.Exec(
			"INSERT INTO profile_technology (profile_id, technology_id) VALUES ($1, $2)",
			profileID,
			technologyID,
		); err != nil {
			return fmt.Errorf("failed to add technology %d to profile %d: %w", technologyID, profileID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit profile technologies transaction: %w", err)
	}

	return nil
}

// GetTechnologyIDs возвращает ID технологий, привязанных к профилю
func (r *ProfileRepo) GetTechnologyIDs(profileID int64) ([]int64, error) {
	const query = `
		SELECT technology_id
		FROM profile_technology
		WHERE profile_id = $1
		ORDER BY technology_id
	`
	rows, err := r.db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to query profile technology ids: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan profile technology id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("profile technology id rows error: %w", err)
	}
	return ids, nil
}
