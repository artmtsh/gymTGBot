package gym

import "fmt"

// Получить или создать UserState для пользователя
func (a *App) getOrCreateUserState(userID int64) *UserState {
	// реализацию напишешь сам
	if _, ok := a.userStates[userID]; !ok {
		a.userStates[userID].State = StateIdle
	}
	return a.userStates[userID]
}

// HandleText обрабатывает любое текстовое сообщение пользователя
// и возвращает текст ответа, который нужно отправить.
func (a *App) HandleText(userID int64, text string) string {
	// внутри:
	//  - распознать команды (/start, /start_workout, /finish_workout, /history, /cancel, /new_exercise)
	//  - посмотреть состояние userState.State
	//  - вызвать нужные методы (StartWorkout, AddExercise, AddSet, FinishWorkout и т.п.)
	switch text {
	case "/start":
		return "Привет я умею выполнять следующие команды - /start, /start_workout, /finish_workout, /history, /cancel, /new_exercise"
	case "start_workout":
		{
			fmt.Println("Введи название тренировки")
			var name string
			//TODO: написать получение названия
			_, err := a.StartWorkout(userID, name)
			if err != nil {
				return "возникла непредвиденная ошибка"
			}
			return "Тренировка " + name + " начата"
		}
	case "finish_workout":
		_, err := a.FinishWorkout(userID)
		if err != nil {
			return "возникла непредвиденная ошибка"
		}
		return "тренировка закончена"
	case "/history":
		var limit int
		//TODO: написать получение количества
		a.GetLastWorkouts(userID, limit)
	case "cancel":
		a.CancelCurrentExercise(userID)
		return "текущее упражнение удалено"
	case "new_exercise":
		var name string
		//TODO: написать получение названия
		_, err := a.AddExercise(userID, name)
		if err != nil {
			return "возникла непредвиденная ошибка"
		}
		return "упражнение " + name + " добавлено"
	default:
		return "неизвестная команда"

	}
	return "что то пошло не так"
}

// Начать тренировку (когда уже знаем название)
func (a *App) StartWorkout(userID int64, name string) (*Workout, error) {
	// создать Workout, сохранить, обновить UserState
	return nil, nil
}

// Добавить упражнение в текущую тренировку пользователя
func (a *App) AddExercise(userID int64, name string) (*Exercise, error) {
	// использовать CurrentWorkoutID, создать Exercise, сохранить, обновить CurrentExerciseID
	return nil, nil
}
