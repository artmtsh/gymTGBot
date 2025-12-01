package gym

import (
	"errors"
	"sync"
	"time"
)

type App struct {
	mu sync.Mutex // на будущее, когда появится параллелизм

	// автоинкрементные ID
	nextWorkoutID  int
	nextExerciseID int
	nextSetID      int

	// данные
	workouts  map[int]*Workout
	exercises map[int]*Exercise
	sets      map[int]*Set

	// состояние пользователей
	userStates map[int64]*UserState
}

func NewApp() *App {
	return &App{
		nextWorkoutID:  1,
		nextExerciseID: 1,
		nextSetID:      1,
		workouts:       make(map[int]*Workout),
		exercises:      make(map[int]*Exercise),
		sets:           make(map[int]*Set),
		userStates:     make(map[int64]*UserState),
	}
}

// Добавить подход к текущему упражнению пользователя
func (a *App) AddSet(userID int64, weight float64, reps int) (*Set, error) {
	// использовать CurrentExerciseID, создать Set, сохранить
	if weight <= 0 {
		return nil, errors.New(NegativeWeight)
	}
	if reps <= 0 {
		return nil, errors.New(NegativeReps)
	}

	var newSet Set
	if a.userStates[userID].State == StateAwaitingSet || a.userStates[userID].CurrentExerciseID != nil {
		newSet = Set{a.nextSetID,
			*a.userStates[userID].CurrentExerciseID,
			weight,
			reps}
		a.sets[a.nextSetID] = &newSet
		a.nextSetID++
	} else {
		return nil, errors.New(NoActiveExercise)
	}
	return &newSet, nil
}

// Завершить текущую тренировку
func (a *App) FinishWorkout(userID int64) (*Workout, error) {
	// проставить FinishedAt, сбросить CurrentWorkoutID / CurrentExerciseID, State = Idle
	var newWorkout Workout
	if a.userStates[userID].State == StateAwaitingExerciseName || a.userStates[userID].State == StateAwaitingSet {
		*a.workouts[*a.userStates[userID].CurrentWorkoutID].FinishedAt = time.Now()
		newWorkout = *a.workouts[*a.userStates[userID].CurrentWorkoutID]
		*a.userStates[userID].CurrentExerciseID = a.nextExerciseID
		*a.userStates[userID].CurrentWorkoutID = a.nextWorkoutID
		a.userStates[userID].State = StateIdle
	} else {
		return nil, errors.New(NoActiveWorkout)
	}
	return &newWorkout, nil
}

func (a *App) CancelCurrentExercise(userID int64) {
	newExercises := make(map[int]*Exercise)
	for k, v := range a.exercises {
		if k != *a.userStates[userID].CurrentExerciseID {
			newExercises[k] = v
		}
	}
	a.exercises = newExercises
}

// История тренировок пользователя (последние N)
func (a *App) GetLastWorkouts(userID int64, limit int) ([]*Workout, error) {
	var lastWorkouts []*Workout
	if *a.userStates[userID].CurrentWorkoutID == 0 {
		return nil, errors.New(ZeroWorkouts)
	} else {
		for i := a.userStates[userID].CurrentWorkoutID; *i < limit; *i-- {
			lastWorkouts = append(lastWorkouts, a.workouts[*i])
		}
	}
	return nil, nil
}
