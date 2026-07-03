package i18n

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// LocalizedText хранит переводы вида {"ru":"...","en":"..."}.
type LocalizedText map[string]string

// LocalizedStringList хранит списки строк по языкам.
type LocalizedStringList map[string][]string

func NewLocalizedText(ru, en string) LocalizedText {
	return LocalizedText{
		string(LocaleRU): strings.TrimSpace(ru),
		string(LocaleEN): strings.TrimSpace(en),
	}
}

func FromLegacyText(value string) LocalizedText {
	return NewLocalizedText(value, "")
}

func (lt LocalizedText) Get(locale Locale) string {
	if lt == nil {
		return ""
	}
	return strings.TrimSpace(lt[string(locale)])
}

func (lt LocalizedText) Set(locale Locale, value string) {
	if lt == nil {
		return
	}
	lt[string(locale)] = strings.TrimSpace(value)
}

func (lt LocalizedText) Resolve(locale Locale) string {
	if text := lt.Get(locale); text != "" {
		return text
	}
	if locale != LocaleRU {
		if text := lt.Get(LocaleRU); text != "" {
			return text
		}
	}
	return lt.Get(LocaleEN)
}

func (lt LocalizedText) IsEmpty() bool {
	return strings.TrimSpace(lt.Get(LocaleRU)) == "" && strings.TrimSpace(lt.Get(LocaleEN)) == ""
}

func (lt LocalizedText) MarshalJSON() ([]byte, error) {
	if lt == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string]string(lt))
}

func (lt *LocalizedText) UnmarshalJSON(data []byte) error {
	if lt == nil {
		return fmt.Errorf("LocalizedText: nil receiver")
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*lt = LocalizedText(raw)
	return nil
}

func (lt *LocalizedText) Scan(value interface{}) error {
	if value == nil {
		*lt = LocalizedText{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return lt.scanBytes(v)
	case string:
		return lt.scanBytes([]byte(v))
	default:
		return fmt.Errorf("LocalizedText: unsupported type %T", value)
	}
}

func (lt *LocalizedText) scanBytes(data []byte) error {
	if len(data) == 0 {
		*lt = LocalizedText{}
		return nil
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*lt = LocalizedText(raw)
	return nil
}

func (lt LocalizedText) Value() (driver.Value, error) {
	if lt == nil {
		return []byte("{}"), nil
	}
	data, err := json.Marshal(map[string]string(lt))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewLocalizedStringList(ru, en []string) LocalizedStringList {
	return LocalizedStringList{
		string(LocaleRU): normalizeLines(ru),
		string(LocaleEN): normalizeLines(en),
	}
}

func FromLegacyStringList(values []string) LocalizedStringList {
	return NewLocalizedStringList(values, nil)
}

func (ll LocalizedStringList) Get(locale Locale) []string {
	if ll == nil {
		return nil
	}
	return append([]string(nil), ll[string(locale)]...)
}

func (ll LocalizedStringList) Set(locale Locale, values []string) {
	if ll == nil {
		return
	}
	ll[string(locale)] = normalizeLines(values)
}

func (ll LocalizedStringList) Resolve(locale Locale) []string {
	if values := ll.Get(locale); len(values) > 0 {
		return values
	}
	if locale != LocaleRU {
		if values := ll.Get(LocaleRU); len(values) > 0 {
			return values
		}
	}
	return ll.Get(LocaleEN)
}

func (ll LocalizedStringList) MarshalJSON() ([]byte, error) {
	if ll == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string][]string(ll))
}

func (ll *LocalizedStringList) UnmarshalJSON(data []byte) error {
	if ll == nil {
		return fmt.Errorf("LocalizedStringList: nil receiver")
	}
	var raw map[string][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*ll = LocalizedStringList(raw)
	return nil
}

func (ll *LocalizedStringList) Scan(value interface{}) error {
	if value == nil {
		*ll = LocalizedStringList{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return ll.scanBytes(v)
	case string:
		return ll.scanBytes([]byte(v))
	default:
		return fmt.Errorf("LocalizedStringList: unsupported type %T", value)
	}
}

func (ll *LocalizedStringList) scanBytes(data []byte) error {
	if len(data) == 0 {
		*ll = LocalizedStringList{}
		return nil
	}
	var raw map[string][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*ll = LocalizedStringList(raw)
	return nil
}

func (ll LocalizedStringList) Value() (driver.Value, error) {
	if ll == nil {
		return []byte("{}"), nil
	}
	data, err := json.Marshal(map[string][]string(ll))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NullableLocalizedText nullable JSONB-поле для опциональных переводов.
type NullableLocalizedText struct {
	Text  LocalizedText
	Valid bool
}

func NullableFromLegacy(value *string) NullableLocalizedText {
	if value == nil || strings.TrimSpace(*value) == "" {
		return NullableLocalizedText{}
	}
	return NullableLocalizedText{
		Text:  FromLegacyText(*value),
		Valid: true,
	}
}

func (n NullableLocalizedText) Resolve(locale Locale) *string {
	if !n.Valid {
		return nil
	}
	text := n.Text.Resolve(locale)
	if text == "" {
		return nil
	}
	return &text
}

func (n NullableLocalizedText) Get(locale Locale) string {
	if !n.Valid {
		return ""
	}
	return n.Text.Get(locale)
}

func (n *NullableLocalizedText) Scan(value interface{}) error {
	if value == nil {
		*n = NullableLocalizedText{}
		return nil
	}
	var text LocalizedText
	if err := text.Scan(value); err != nil {
		return err
	}
	if text.IsEmpty() {
		*n = NullableLocalizedText{}
		return nil
	}
	*n = NullableLocalizedText{Text: text, Valid: true}
	return nil
}

func (n NullableLocalizedText) Value() (driver.Value, error) {
	if !n.Valid || n.Text.IsEmpty() {
		return nil, nil
	}
	return n.Text.Value()
}

func normalizeLines(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func ParseLinesInput(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return normalizeLines(strings.Split(value, "\n"))
}
