package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/exercise/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {

	var duration time.Duration
	var err error

	partsData := strings.Split(data, ",")
	if len(partsData) != 2 {
		return 0, 0, fmt.Errorf("error partsData need 2 parts, got %d", len(partsData))
	}

	stepsStr := strings.TrimSpace(partsData[0])
	durationStr := strings.TrimSpace(partsData[1])

	var steps int
	steps, err = strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("error convert stepStr to integer")
	}

	duration, err = time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("error convert durationStr to time.Duration")
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Errorf("error DayActionInfo, failed parsePackage(data)")
		return ""
	}

	if steps <= 0 {
		return fmt.Sprintf("steps <0 err: %s", err.Error())
	}
	distance := float64(steps) * StepLength
	kmDistance := distance / 1000
	kalories := spentcalories.WalkingSpentCalories(steps, weight, height, duration) // (steps int, weight, height float64, duration time.Duration) float64

	return fmt.Sprintf(" Колличество шагов: %d\n Дистанция составила %.2fкм.\n Вы сожгли %.2f ккал.", steps, kmDistance, kalories)
}
