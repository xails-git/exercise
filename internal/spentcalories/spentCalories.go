package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {

	var steps int
	var species string
	var duration time.Duration
	var err error

	partsData := strings.Split(data, ",")
	if len(partsData) != 3 {
		return 0, "", 0, fmt.Errorf("error partsData, need len 3, have %d", len(partsData))
	}

	stepsStr := strings.TrimSpace(partsData[0])
	speciesStr := strings.TrimSpace(partsData[1]) // удаляем пробелы
	durationStr := strings.TrimSpace(partsData[2])

	species = speciesStr // сразу поместим в нужную переменную
	steps, err = strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error convert stepsStr to integer")
	}

	duration, err = time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error convert durationStr to time.duration")
	}
	return steps, species, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	result := (float64(steps) * lenStep) / mInKm
	return result
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 || steps == 0 {
		return 0
	}
	distance := distance(steps)
	hourDuration := duration.Hours() // переводим в часы float64

	avgHoursPerUnit := distance / hourDuration // вычисляем ср. скорость
	return avgHoursPerUnit

}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {

	var resultCalories float64
	var resultDistance float64
	var resultAvgSpeed float64

	steps, species, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Data parsing error: %s", err.Error())
	}
	durationHours := duration.Hours()

	switch species {
	case "Бег":
		resultCalories = RunningSpentCalories(steps, weight, duration)
		resultDistance = distance(steps)
		resultAvgSpeed = meanSpeed(steps, duration)

	case "Ходьба":
		resultDistance = distance(steps)
		resultAvgSpeed = meanSpeed(steps, duration)
		resultCalories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		species,        // тип тренировки (string)
		durationHours,  // длительность в часах (float64)
		resultDistance, // дистанция (float64)
		resultAvgSpeed, // скорость (float64)
		resultCalories, // калории (float64)
	)

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	avgHoursPerUnit := meanSpeed(steps, duration)
	if avgHoursPerUnit <= 0 || weight <= 0 {
		return 0
	}
	calories := ((runningCaloriesMeanSpeedMultiplier * avgHoursPerUnit) - runningCaloriesMeanSpeedShift) * weight
	return calories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	if weight <= 0 || height <= 0 || duration <= 0 || steps <= 0 {
		return 0
	}
	avgHoursPerUnit := meanSpeed(steps, duration)
	if avgHoursPerUnit <= 0 {
		return 0
	}
	durationHours := duration.Hours() // превращаеи duration в часы float64

	calories := ((walkingCaloriesWeightMultiplier * weight) + (avgHoursPerUnit*avgHoursPerUnit/height)*walkingSpeedHeightMultiplier) * durationHours * minInH
	return calories
}
