package i18n

// PDFSectionLabels возвращает локализованные заголовки секций PDF.
func PDFSectionLabels(locale Locale) map[string]string {
	if locale == LocaleEN {
		return map[string]string{
			"about":        "About me",
			"work":         "Work experience",
			"technologies": "Technologies",
			"education":    "Education",
			"contacts":     "Contacts",
			"present":      "present",
		}
	}
	return map[string]string{
		"about":        "О себе",
		"work":         "Опыт работы",
		"technologies": "Технологии",
		"education":    "Образование",
		"contacts":     "Контакты",
		"present":      "по наст. время",
	}
}
