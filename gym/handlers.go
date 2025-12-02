package gym

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Получить или создать UserState для пользователя
func (a *App) getOrCreateUserState(userID int64) *UserState {
	if st, ok := a.userStates[userID]; ok {
		return st
	}

	st := &UserState{
		UserID: userID,
		State:  StateIdle,
	}
	a.userStates[userID] = st
	return st
}

// HandleText обрабатывает любое текстовое сообщение пользователя
// и возвращает текст ответа, который нужно отправить.
func (a *App) HandleText(userID int64, text string) string {
	// внутри:
	//  - распознать команды (/start, /start_workout, /finish_workout, /history, /cancel, /new_exercise)
	//  - посмотреть состояние userState.State
	//  - вызвать нужные методы (StartWorkout, AddExercise, AddSet, FinishWorkout и т.п.)
	user := a.getOrCreateUserState(userID)
	arrText := strings.Split(text, " ")
	if runes := arrText[0]; runes[0] == '/' {
		switch text {
		case "/start":
			return "Привет я умею выполнять следующие команды - /start, /start_workout, /finish_workout, /history, /cancel, /new_exercise"
		case "/start_workout":
			{
				user.State = StateAwaitingWorkoutName
				return "Введите название тренировки"
			}
		case "/finish_workout":
			_, err := a.FinishWorkout(userID)
			if err != nil {
				return "возникла непредвиденная ошибка"
			}
			user.State = StateIdle
			return "тренировка завершена"
		case "/history":
			user.State = StateAwaitingLimitForHistory
			return "введите количество тренировок о которых хотите узнать"
		case "cancel":
			err := a.CancelCurrentExercise(userID)
			user.State = StateAwaitingExerciseName
			if err != nil {
				return ""
			}
			return "текущее упражнение удалено"
		case "/new_exercise":
			user.State = StateAwaitingExerciseName
			return "Введите название упражнения"
		default:
			return "неизвестная команда"
		}
	} else {
		switch user.State {
		case StateIdle:
			return UnknownAction
		case StateAwaitingLimitForHistory:
			if len(arrText) != 1 {
				return MoreThanOneValue
			}
			n, err := strconv.Atoi(text)
			if err != nil {
				return NotANumber
			}
			if n < 0 {
				return NegativeNumber
			}
			workouts, err := a.GetLastWorkouts(userID, n)
			if err != nil {
				return ""
			}
			return a.workoutsToString(workouts)
		case StateAwaitingExerciseName:
			_, err := a.AddExercise(userID, text)
			if err != nil {
				return ""
			}
			return "упражнение " + text + " добавлено"
		case StateAwaitingSet:
			if len(arrText) != 2 {
				return IllegalAmountOfArgs
			}
			weight, _ := strconv.ParseFloat(arrText[0], 64)
			reps, _ := strconv.Atoi(arrText[0])
			_, err := a.AddSet(userID, weight, reps)
			if err != nil {
				return ""
			}
			return "введите следующий подход или закончите упражнение"
		}
	}

	return "что то пошло не так"
}

// Начать тренировку (когда уже знаем название)
func (a *App) StartWorkout(userID int64, name string) (*Workout, error) {
	// создать Workout, сохранить, обновить UserState
	if strings.TrimSpace(name) == "" {
		return nil, errors.New(EmptyName)
	}
	st := a.getOrCreateUserState(userID)
	if st.State != StateAwaitingWorkoutName {
		return nil, errors.New(ImpossibleToStartNewWorkout)
	}
	st.State = StateAwaitingExerciseName
	newWorkout := Workout{ID: a.nextWorkoutID,
		UserID:     userID,
		Name:       name,
		StartedAt:  time.Now(),
		FinishedAt: nil,
	}
	a.workouts[a.nextWorkoutID] = &newWorkout
	a.nextWorkoutID++
	return &newWorkout, nil
}

// Добавить упражнение в текущую тренировку пользователя
func (a *App) AddExercise(userID int64, name string) (*Exercise, error) {
	// использовать CurrentWorkoutID, создать Exercise, сохранить, обновить CurrentExerciseID
	st := a.getOrCreateUserState(userID)
	if st.State != StateAwaitingExerciseName && st.State != StateAwaitingSet {
		return nil, errors.New(ImpossibleToStartNewExercise)
	}
	if st.CurrentWorkoutID == nil {
		return nil, errors.New(NoActiveWorkout)
	}
	newExercise := Exercise{ID: a.nextExerciseID,
		WorkoutID: *st.CurrentWorkoutID,
		Name:      name}
	a.exercises[*st.CurrentExerciseID] = &newExercise
	*st.CurrentExerciseID++
	a.nextExerciseID++
	return &newExercise, nil
}
