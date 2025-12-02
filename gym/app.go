package gym

import (
	"errors"
	"strconv"
	"strings"
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
	st := a.getOrCreateUserState(userID)
	var newSet Set
	if st.State == StateAwaitingSet || st.CurrentExerciseID != nil {
		newSet = Set{a.nextSetID,
			*st.CurrentExerciseID,
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
	var tempWorkout Workout
	st := a.getOrCreateUserState(userID)
	if st.CurrentWorkoutID != nil || st.State != StateIdle {
		w := a.workouts[*st.CurrentWorkoutID]
		if w.FinishedAt == nil {
			t := time.Now()
			w.FinishedAt = &t
		}
	} else {
		return nil, errors.New(NoActiveWorkout)
	}
	st.CurrentWorkoutID = nil
	st.CurrentExerciseID = nil
	st.State = StateIdle
	return &tempWorkout, nil
}

func (a *App) CancelCurrentExercise(userID int64) error {
	st := a.getOrCreateUserState(userID)
	if st.State != StateAwaitingSet {
		return errors.New(NoActiveExercise)
	}
	newSets := make(map[int]*Set)
	for k, v := range a.sets {
		if k != *a.userStates[userID].CurrentExerciseID {
			newSets[k] = v
		}
	}
	st.State = StateAwaitingExerciseName
	a.sets = newSets
	return nil
}

// История тренировок пользователя (последние N)
func (a *App) GetLastWorkouts(userID int64, limit int) ([]*Workout, error) {
	var lastWorkouts []*Workout
	for _, v := range a.workouts {
		if v.UserID == userID && v.FinishedAt != nil {
			lastWorkouts = append(lastWorkouts, v)
		}
	}
	if len(lastWorkouts) == 0 {
		return nil, errors.New(ZeroWorkouts)
	}
	return lastWorkouts[:limit], nil
}

func (a *App) workoutsToString(workouts []*Workout) string {
	var resString strings.Builder
	for _, workout := range workouts {
		resString.WriteString("Тренировка ")
		resString.WriteString(workout.Name)
		resString.WriteByte('\n')
		for _, exercise := range a.exercises {
			if exercise.WorkoutID == workout.ID {
				resString.WriteString("\t Упражнение ")
				resString.WriteString(exercise.Name)
				resString.WriteByte(':')
				resString.WriteByte('\n')
				for _, set := range a.sets {
					if set.ExerciseID == exercise.ID {
						resString.WriteString("/t /t Подход: \n")
						resString.WriteString("\t \t \t")
						resString.WriteString(strconv.Itoa(set.Reps))
						resString.WriteString(" ")
						resString.WriteString(strconv.FormatFloat(set.Weight, 'f', -1, 64))
						resString.WriteByte('\n')
					}
				}
			}
		}
	}
	return resString.String()
}
