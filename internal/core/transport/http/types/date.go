package core_http_types

import (
	"bytes"
	"time"
)

// Date — это алиас над time.Time, который живет только в HTTP-слое
type Date time.Time

// Указываем нужный формат (в Go 2006-01-02 — это эталонная дата)
const dateLayout = time.DateOnly

// UnmarshalJSON учит декодер правильно читать нашу строку
func (d *Date) UnmarshalJSON(b []byte) error {
	// JSON присылает строку в кавычках: "2026-05-05". Отрезаем их.
	str := string(bytes.Trim(b, `"`))
	
	// Защита от пустоты
	if str == "null" || str == "" {
		return nil
	}

	// Парсим короткую строку в полноценный объект time.Time
	t, err := time.Parse(dateLayout, str)
	if err != nil {
		return err
	}

	// Сохраняем результат
	*d = Date(t)
	return nil
}

// MarshalJSON учит декодер правильно отдавать дату обратно фронтенду
func (d Date) MarshalJSON() ([]byte, error) {
	// Превращаем Date обратно в time.Time и форматируем
	formatted := time.Time(d).Format(dateLayout)
	// Оборачиваем в кавычки для валидного JSON
	return []byte(`"` + formatted + `"`), nil
}
