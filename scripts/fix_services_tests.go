package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	root := filepath.Join("internal", "services")
	files := []string{"education_test.go", "technology_test.go", "workhist_test.go", "tag_test.go"}

	for _, name := range files {
		path := filepath.Join(root, name)
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		text := string(content)

		if !strings.Contains(text, "pkg/i18n") {
			text = strings.Replace(text,
				`entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"`,
				`entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"`, 1)
		}

		rePgName := regexp.MustCompile(`Name:\s+pgtype\.Text\{String: "([^"]+)", Valid: true\}`)
		text = rePgName.ReplaceAllString(text, `Name: newNullableLocalizedText("$1")`)

		reDesc := regexp.MustCompile(`Description: pgtype\.Text\{String: "([^"]+)", Valid: true\}`)
		text = reDesc.ReplaceAllString(text, `Description: newNullableLocalizedText("$1")`)

		reLogo := regexp.MustCompile(`LogoUrl:\s+pgtype\.Text\{String: "([^"]+)", Valid: true\}`)
		text = reLogo.ReplaceAllString(text, `LogoUrl: newPgText("$1")`)

		reCourse := regexp.MustCompile(`Course:\s+"([^"]*)"`)
		text = reCourse.ReplaceAllString(text, `Course: newLocalizedText("$1")`)

		reOrg := regexp.MustCompile(`Organization:\s+"([^"]*)"`)
		text = reOrg.ReplaceAllString(text, `Organization: newLocalizedText("$1")`)

		reTagName := regexp.MustCompile(`(models\.Tag\{[^}]*|\t\t\t\tName:\s+)("Backend"|"Frontend"|"Duplicate"|"TestGet"|"Original"|"Updated"|"NonExistent"|"ToDelete"|"Tag1"|"Tag2"|"Tag3"|"Alpha"|"Beta"|"Gamma"|"Delta"|"Epsilon"|"Charlie"|"Bravo")`)
		text = reTagName.ReplaceAllStringFunc(text, func(s string) string {
			if strings.Contains(s, "models.Tag") {
				return strings.Replace(s, `Name: "`, `Name: newLocalizedText("`, 1)
			}
			return s
		})

		// Tag Name fields in struct literals
		text = regexp.MustCompile(`Name:\s+"([^"]+)"`).ReplaceAllStringFunc(text, func(s string) string {
			if strings.Contains(s, `name: "`) {
				return s
			}
			sub := regexp.MustCompile(`Name:\s+"([^"]+)"`).FindStringSubmatch(s)
			if len(sub) != 2 {
				return s
			}
			val := sub[1]
			if val == "Backend" || val == "Frontend" || val == "Go" || val == "Docker" || val == "React" ||
				strings.Contains(val, "Company") || strings.Contains(val, "Яндекс") || strings.Contains(val, "Tinkoff") {
				return `Name: newLocalizedText("` + val + `")`
			}
			return s
		})

		reAbout := regexp.MustCompile(`About:\s+"([^"]+)"`)
		text = reAbout.ReplaceAllString(text, `About: newLocalizedText("$1")`)

		reWhat := regexp.MustCompile(`WhatIDid:\s+\[\]string\{([^}]*)\}`)
		text = reWhat.ReplaceAllString(text, `WhatIDid: newLocalizedStringList($1)`)

		reProjects := regexp.MustCompile(`Projects:\s+\[\]string\{([^}]*)\}`)
		text = reProjects.ReplaceAllString(text, `Projects: newLocalizedStringList($1)`)

		text = strings.ReplaceAll(text,
			`if result.Course != tt.mockEdu.Course {`,
			`if result.Course.Get(i18n.LocaleRU) != tt.mockEdu.Course.Get(i18n.LocaleRU) {`)
		text = strings.ReplaceAll(text,
			`t.Errorf("Ожидался Course = %s, получили %s", tt.mockEdu.Course, result.Course)`,
			`t.Errorf("Ожидался Course = %s, получили %s", tt.mockEdu.Course.Get(i18n.LocaleRU), result.Course.Get(i18n.LocaleRU))`)

		text = strings.ReplaceAll(text,
			`if len(result.WhatIDid) != len(tt.mockWH.WhatIDid) {`,
			`if len(result.WhatIDid.Get(i18n.LocaleRU)) != len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)) {`)
		text = strings.ReplaceAll(text,
			`len(tt.mockWH.WhatIDid), len(result.WhatIDid)`,
			`len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)), len(result.WhatIDid.Get(i18n.LocaleRU))`)

		text = strings.ReplaceAll(text,
			`func (m *MockTechRepo) GetWithTags(id int64) (dto.TechnologyWithTagsDTO, error) {`,
			`func (m *MockTechRepo) GetWithTags(id int64, locale i18n.Locale) (dto.TechnologyWithTagsDTO, error) {`)
		text = strings.ReplaceAll(text,
			`func (m *MockTechRepo) ListWithTags(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO], error) {`,
			`func (m *MockTechRepo) ListWithTags(req entityreqdecorator.PagebleRq, locale i18n.Locale) (entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO], error) {`)

		text = strings.ReplaceAll(text,
			`func (m *MockWorkHistoryRepo) GetWithTechnologies(id int64) (dto.WorkHistoryWithTechnologiesDTO, error) {`,
			`func (m *MockWorkHistoryRepo) GetWithTechnologies(id int64, locale i18n.Locale) (dto.WorkHistoryWithTechnologiesDTO, error) {`)
		text = strings.ReplaceAll(text,
			`func (m *MockWorkHistoryRepo) ListWithTechnologies(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO], error) {`,
			`func (m *MockWorkHistoryRepo) ListWithTechnologies(req entityreqdecorator.PagebleRq, locale i18n.Locale) (entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO], error) {`)

		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			panic(err)
		}
	}
}
