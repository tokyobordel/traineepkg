package trace

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	// Алфавит для генерации (без похожих символов: 0, O, 1, I, l)
	alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

	// Длина компонентов
	separatorLength = 1  // символ "_"
	randomLength    = 10 // символов для случайной части
	// timestampLength теперь переменная - зависит от длины base64 timestamp
)

// GenerateSpreadId генерирует уникальный spread ID с переменной длиной
// Формат: timestamp в base64 + "_" + 12 случайных символов
// Пример: MTcwMzQyNTY3ODkwMQ + _ + 8G9H2J3K4L5M = MTcwMzQyNTY3ODkwMQ_8G9H2J3K4L5M
func GenerateSpreadId() string {
	// Получаем timestamp в наносекундах для максимальной точности
	now := time.Now().UnixNano()

	// Конвертируем timestamp в байты
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(now))

	// Кодируем в base64 (убираем padding для компактности)
	timestampStr := strings.TrimRight(base64.StdEncoding.EncodeToString(timestampBytes), "=")

	// Генерируем криптографически стойкую случайную часть
	randomPart := generateRandomString(randomLength)

	// Объединяем timestamp, подчеркивание и случайную часть
	return timestampStr + "_" + randomPart
}

// generateRandomString генерирует криптографически стойкую случайную строку
func generateRandomString(length int) string {
	bytes := make([]byte, length)

	// Используем crypto/rand для максимальной энтропии
	if _, err := rand.Read(bytes); err != nil {
		// Fallback на псевдослучайный генератор в случае ошибки
		return generateFallbackRandomString(length)
	}

	// Конвертируем байты в символы из алфавита
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphabet[int(bytes[i])%len(alphabet)]
	}

	return string(result)
}

// generateFallbackRandomString - запасной генератор на основе времени
func generateFallbackRandomString(length int) string {
	// Используем микросекунды как seed для псевдослучайности
	seed := time.Now().UnixMicro()

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		// Простой LCG генератор
		seed = (seed*1103515245 + 12345) & 0x7fffffff
		result[i] = alphabet[seed%int64(len(alphabet))]
	}

	return string(result)
}

// ParseSpreadId извлекает информацию из spread ID (для отладки)
// Возвращает timestamp в наносекундах
func ParseSpreadId(spreadId string) (timestamp int64, randomPart string, err error) {
	// Находим позицию подчеркивания
	separatorIndex := strings.Index(spreadId, "_")
	if separatorIndex == -1 {
		return 0, "", fmt.Errorf("invalid spread ID format: separator '_' not found")
	}

	// Проверяем, что после подчеркивания есть случайная часть нужной длины
	if len(spreadId)-separatorIndex-1 != randomLength {
		return 0, "", fmt.Errorf("invalid spread ID format: random part should be %d characters, got %d", randomLength, len(spreadId)-separatorIndex-1)
	}

	timestampStr := spreadId[:separatorIndex]
	randomPart = spreadId[separatorIndex+1:] // +1 чтобы пропустить подчеркивание

	// Добавляем padding для base64 если нужно
	for len(timestampStr)%4 != 0 {
		timestampStr += "="
	}

	// Декодируем base64 обратно в байты
	timestampBytes, err := base64.StdEncoding.DecodeString(timestampStr)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode base64 timestamp: %w", err)
	}

	// Проверяем размер
	if len(timestampBytes) != 8 {
		return 0, "", fmt.Errorf("invalid timestamp bytes length: expected 8, got %d", len(timestampBytes))
	}

	// Конвертируем байты обратно в int64
	timestamp = int64(binary.BigEndian.Uint64(timestampBytes))

	return timestamp, randomPart, nil
}
