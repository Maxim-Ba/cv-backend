package services

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/repository"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

//go:embed fonts/DejaVuSans.ttf
var dejaVuRegular []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var dejaVuBold []byte

const (
	pdfMargin      = 15.0
	pdfContentW    = 180.0
	pdfLineH       = 5.5
	pdfHeaderBgR   = 41
	pdfHeaderBgG   = 65
	pdfHeaderBgB   = 122
)

// ProfileGetter интерфейс для получения профиля
type ProfileGetter interface {
	Get() (repository.Profile, error)
}

// AboutMeGetter интерфейс для получения секции «О себе»
type AboutMeGetter interface {
	GetAboutMe() (dto.AboutMeDTO, error)
}

// PDFService генерирует PDF-резюме из данных БД
type PDFService struct {
	profile ProfileGetter
	aboutMe AboutMeGetter
	wh      *WorkHistoryService
	tech    *TechService
	edu     *EducationService
}

// NewPDFService создает новый экземпляр сервиса генерации PDF
func NewPDFService(
	profile ProfileGetter,
	aboutMe AboutMeGetter,
	wh *WorkHistoryService,
	tech *TechService,
	edu *EducationService,
) *PDFService {
	return &PDFService{
		profile: profile,
		aboutMe: aboutMe,
		wh:      wh,
		tech:    tech,
		edu:     edu,
	}
}

// GenerateCV собирает данные из БД и генерирует PDF с CV
func (s *PDFService) GenerateCV() ([]byte, error) {
	profile, err := s.profile.Get()
	if err != nil {
		return nil, fmt.Errorf("pdf: get profile: %w", err)
	}

	allPage := entityreqdecorator.PagebleRq{Page: 1, Size: 200}

	whResult, err := s.wh.ListWithTechnologies(allPage)
	if err != nil {
		return nil, fmt.Errorf("pdf: get work history: %w", err)
	}

	techResult, err := s.tech.ListWithTags(allPage)
	if err != nil {
		return nil, fmt.Errorf("pdf: get technologies: %w", err)
	}

	eduResult, err := s.edu.List(allPage)
	if err != nil {
		return nil, fmt.Errorf("pdf: get education: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("DejaVu", "", dejaVuRegular)
	pdf.AddUTF8FontFromBytes("DejaVu", "B", dejaVuBold)
	pdf.SetMargins(pdfMargin, pdfMargin, pdfMargin)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	pdfRenderHeader(pdf, profile)

	aboutMe, err := s.aboutMe.GetAboutMe()
	if err != nil {
		return nil, fmt.Errorf("pdf: get about me: %w", err)
	}
	if pdfAboutMeHasContent(aboutMe) {
		pdfSectionTitle(pdf, "О СЕБЕ")
		pdfAboutMe(pdf, aboutMe)
	}

	if len(whResult.Content) > 0 {
		pdfSectionTitle(pdf, "ОПЫТ РАБОТЫ")
		for _, w := range whResult.Content {
			pdfWorkItem(pdf, w)
		}
	}

	if len(techResult.Content) > 0 {
		pdfSectionTitle(pdf, "ТЕХНОЛОГИИ")
		pdfTechnologies(pdf, techResult.Content)
	}

	if len(eduResult.Content) > 0 {
		pdfSectionTitle(pdf, "ОБРАЗОВАНИЕ")
		for _, e := range eduResult.Content {
			pdfEducationItem(pdf, e)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf: output: %w", err)
	}
	return buf.Bytes(), nil
}

// pdfRenderHeader рисует шапку с именем, должностью и контактами
func pdfRenderHeader(pdf *fpdf.Fpdf, p repository.Profile) {
	pdf.SetFillColor(pdfHeaderBgR, pdfHeaderBgG, pdfHeaderBgB)
	pdf.Rect(0, 0, 210, 44, "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(pdfMargin, 8)
	pdf.SetFont("DejaVu", "B", 22)
	pdf.CellFormat(pdfContentW, 10, p.FullName, "", 1, "L", false, 0, "")

	pdf.SetFont("DejaVu", "", 13)
	pdf.SetX(pdfMargin)
	pdf.CellFormat(pdfContentW, 7, p.Title, "", 1, "L", false, 0, "")

	contacts := pdfBuildContacts(p)
	if contacts != "" {
		pdf.SetFont("DejaVu", "", 8)
		pdf.SetX(pdfMargin)
		pdf.CellFormat(pdfContentW, 6, contacts, "", 1, "L", false, 0, "")
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(50)
}

func pdfBuildContacts(p repository.Profile) string {
	sep := "   |   "
	var parts []string
	if p.Email != nil && *p.Email != "" {
		parts = append(parts, *p.Email)
	}
	if p.Phone != nil && *p.Phone != "" {
		parts = append(parts, *p.Phone)
	}
	if p.Telegram != nil && *p.Telegram != "" {
		parts = append(parts, *p.Telegram)
	}
	if p.GitHub != nil && *p.GitHub != "" {
		parts = append(parts, *p.GitHub)
	}
	return strings.Join(parts, sep)
}

func pdfAboutMeHasContent(aboutMe dto.AboutMeDTO) bool {
	if len(aboutMe.BioParagraphs) > 0 {
		return true
	}
	if len(aboutMe.Technologies) > 0 {
		return true
	}
	if aboutMe.Note != nil && strings.TrimSpace(*aboutMe.Note) != "" {
		return true
	}
	if aboutMe.Hobbies != nil && strings.TrimSpace(*aboutMe.Hobbies) != "" {
		return true
	}
	return false
}

// pdfAboutMe рисует секцию «О себе»: биография, бейджи технологий, заметка и хобби
func pdfAboutMe(pdf *fpdf.Fpdf, aboutMe dto.AboutMeDTO) {
	pdf.SetFont("DejaVu", "", 9)
	for _, paragraph := range aboutMe.BioParagraphs {
		if strings.TrimSpace(paragraph) == "" {
			continue
		}
		pdf.SetX(pdfMargin)
		pdf.MultiCell(pdfContentW, pdfLineH, paragraph, "", "L", false)
		pdf.SetY(pdf.GetY() + 1)
	}

	if len(aboutMe.Technologies) > 0 {
		names := make([]string, 0, len(aboutMe.Technologies))
		for _, t := range aboutMe.Technologies {
			names = append(names, t.Title)
		}
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "B", 9)
		pdf.CellFormat(24, pdfLineH, "Стек:", "", 0, "L", false, 0, "")
		pdf.SetFont("DejaVu", "", 9)
		pdf.MultiCell(pdfContentW-24, pdfLineH, strings.Join(names, "  •  "), "", "L", false)
	}

	if aboutMe.Note != nil && strings.TrimSpace(*aboutMe.Note) != "" {
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetTextColor(90, 90, 90)
		pdf.MultiCell(pdfContentW, pdfLineH, *aboutMe.Note, "", "L", false)
		pdf.SetTextColor(0, 0, 0)
	}

	if aboutMe.Hobbies != nil && strings.TrimSpace(*aboutMe.Hobbies) != "" {
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetTextColor(90, 90, 90)
		pdf.MultiCell(pdfContentW, pdfLineH, *aboutMe.Hobbies, "", "L", false)
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetY(pdf.GetY() + 2)
}

// pdfSectionTitle рисует заголовок секции с цветным фоном
func pdfSectionTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetY(pdf.GetY() + 3)
	pdf.SetFillColor(235, 237, 245)
	pdf.SetTextColor(pdfHeaderBgR, pdfHeaderBgG, pdfHeaderBgB)
	pdf.SetFont("DejaVu", "B", 10)
	pdf.SetX(pdfMargin)
	pdf.CellFormat(pdfContentW, 7, title, "", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(pdf.GetY() + 2)
}

// pdfWorkItem рисует одну запись опыта работы
func pdfWorkItem(pdf *fpdf.Fpdf, w dto.WorkHistoryWithTechnologiesDTO) {
	period := pdfFormatPeriod(w.PeriodStart, w.PeriodEnd)

	pdf.SetX(pdfMargin)
	pdf.SetFont("DejaVu", "B", 10)
	pdf.CellFormat(125, 6, w.Name, "", 0, "L", false, 0, "")
	pdf.SetFont("DejaVu", "", 8)
	pdf.SetTextColor(90, 90, 90)
	pdf.CellFormat(55, 6, period, "", 1, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	if w.About != "" {
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "", 9)
		pdf.MultiCell(pdfContentW, pdfLineH, w.About, "", "L", false)
	}

	if len(w.WhatIDid) > 0 {
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "B", 9)
		pdf.CellFormat(pdfContentW, pdfLineH, "Что я делал:", "", 1, "L", false, 0, "")
		pdf.SetFont("DejaVu", "", 9)
		for _, item := range w.WhatIDid {
			if strings.TrimSpace(item) == "" {
				continue
			}
			pdf.SetX(pdfMargin + 4)
			pdf.MultiCell(pdfContentW-4, pdfLineH, "• "+item, "", "L", false)
		}
	}

	if len(w.Technologies) > 0 {
		names := make([]string, 0, len(w.Technologies))
		for _, t := range w.Technologies {
			names = append(names, t.Title)
		}
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "B", 9)
		pdf.CellFormat(24, pdfLineH, "Технологии:", "", 0, "L", false, 0, "")
		pdf.SetFont("DejaVu", "", 9)
		pdf.MultiCell(pdfContentW-24, pdfLineH, strings.Join(names, ", "), "", "L", false)
	}

	pdf.SetDrawColor(210, 210, 220)
	pdf.Line(pdfMargin, pdf.GetY()+1, 210-pdfMargin, pdf.GetY()+1)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetY(pdf.GetY() + 4)
}

// pdfTechnologies рисует технологии, сгруппированные по тегам
func pdfTechnologies(pdf *fpdf.Fpdf, techs []dto.TechnologyWithTagsDTO) {
	tagMap := make(map[string][]string)
	tagOrder := []string{}

	for _, t := range techs {
		if len(t.Tags) == 0 {
			const other = "Другое"
			if _, ok := tagMap[other]; !ok {
				tagOrder = append(tagOrder, other)
			}
			if !pdfContains(tagMap[other], t.Title) {
				tagMap[other] = append(tagMap[other], t.Title)
			}
			continue
		}
		for _, tag := range t.Tags {
			if _, ok := tagMap[tag.Name]; !ok {
				tagOrder = append(tagOrder, tag.Name)
			}
			if !pdfContains(tagMap[tag.Name], t.Title) {
				tagMap[tag.Name] = append(tagMap[tag.Name], t.Title)
			}
		}
	}

	for _, tagName := range tagOrder {
		list := tagMap[tagName]
		pdf.SetX(pdfMargin)
		pdf.SetFont("DejaVu", "B", 9)
		label := "[" + tagName + "]"
		pdf.CellFormat(32, pdfLineH, label, "", 0, "L", false, 0, "")
		pdf.SetFont("DejaVu", "", 9)
		pdf.MultiCell(pdfContentW-32, pdfLineH, strings.Join(list, "  •  "), "", "L", false)
	}
	pdf.SetY(pdf.GetY() + 2)
}

// pdfEducationItem рисует одну запись образования
func pdfEducationItem(pdf *fpdf.Fpdf, e models.Education) {
	nameStr := ""
	if e.Name.Valid {
		nameStr = e.Name.String
	}
	line := e.Course + " — " + e.Organization
	if nameStr != "" {
		line = nameStr + ": " + line
	}

	pdf.SetX(pdfMargin)
	pdf.SetFont("DejaVu", "B", 9)
	pdf.CellFormat(14, pdfLineH, fmt.Sprintf("%d", e.Year), "", 0, "L", false, 0, "")
	pdf.SetFont("DejaVu", "", 9)
	pdf.MultiCell(pdfContentW-14, pdfLineH, line, "", "L", false)
}

// pdfFormatPeriod форматирует диапазон дат для отображения в PDF
func pdfFormatPeriod(start, end *string) string {
	s := "—"
	if start != nil && len(*start) >= 7 {
		s = pdfFormatDate(*start)
	}
	e := "по настоящее время"
	if end != nil && *end != "" {
		e = pdfFormatDate(*end)
	}
	return s + " – " + e
}

var pdfMonthRu = [13]string{
	"", "янв", "фев", "мар", "апр", "май", "июн",
	"июл", "авг", "сен", "окт", "ноя", "дек",
}

// pdfFormatDate преобразует ISO-дату "2021-06-01" в "июн 2021"
func pdfFormatDate(iso string) string {
	if len(iso) < 7 {
		return iso
	}
	var year, month int
	if _, err := fmt.Sscanf(iso[:7], "%d-%d", &year, &month); err != nil || month < 1 || month > 12 {
		return iso[:7]
	}
	return fmt.Sprintf("%s %d", pdfMonthRu[month], year)
}

func pdfContains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
