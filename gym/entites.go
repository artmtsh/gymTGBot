package gym

import "time"

// Тренировка
type Workout struct {
	ID         int        `json:"id"`
	UserID     int64      `json:"user_id"` // Telegram user ID
	Name       string     `json:"name"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"` // nil, если тренировка ещё идёт
}

// Упражнение внутри тренировки
type Exercise struct {
	ID        int    `json:"id"`
	WorkoutID int    `json:"workout_id"`
	Name      string `json:"name"`
}

// Подход
type Set struct {
	ID         int     `json:"id"`
	ExerciseID int     `json:"exercise_id"`
	Weight     float64 `json:"weight"`
	Reps       int     `json:"reps"`
}
