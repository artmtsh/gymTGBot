package gym

// Состояние диалога с пользователем
type DialogState int

const (
	StateIdle                 DialogState = iota
	StateAwaitingWorkoutName              // ждём название тренировки после /start_workout
	StateAwaitingExerciseName             // ждём название упражнения
	StateAwaitingSet                      // ждём "вес повторы" для текущего упражнения
)

// Состояние конкретного пользователя
type UserState struct {
	UserID            int64       // Telegram user ID
	State             DialogState // текущее состояние
	CurrentWorkoutID  *int        // активная тренировка (если есть)
	CurrentExerciseID *int        // текущее упражнение (если есть)
}
