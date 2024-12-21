package db

import (
	"app/log"
	"app/model"
)

func GetSearchData() (*model.DataResponse, error) {
	log.App.Info("Начало получения данных для поиска")

	// Получаем предметы
	var subjects []model.Subject
	log.App.Info("Получение списка предметов")
	if err := App.Find(&subjects).Error; err != nil {
		log.App.Error("Ошибка при получении предметов:", err)
		return nil, err
	}
	log.App.Info("Получено предметов:", len(subjects))

	// Получаем уровни подготовки
	var level_trainings []model.LevelTraining
	log.App.Info("Получение списка уровней подготовки")
	if err := App.Find(&level_trainings).Error; err != nil {
		log.App.Error("Ошибка при получении уровней подготовки:", err)
		return nil, err
	}
	log.App.Info("Получено уровней подготовки:", len(level_trainings))

	// Получаем уровни опыта
	var level_experience []model.TutorExperience
	log.App.Info("Получение списка уровней опыта")
	if err := App.Find(&level_experience).Error; err != nil {
		log.App.Error("Ошибка при получении уровней опыта:", err)
		return nil, err
	}
	log.App.Info("Получено уровней опыта:", len(level_experience))

	// Получаем форматы занятий
	var class_formats []model.ClassFormat
	log.App.Info("Получение списка форматов занятий")
	if err := App.Find(&class_formats).Error; err != nil {
		log.App.Error("Ошибка при получении форматов занятий:", err)
		return nil, err
	}
	log.App.Info("Получено форматов занятий:", len(class_formats))

	log.App.Info("Формирование ответа")
	response := &model.DataResponse{
		WebSocketResponse: model.WebSocketResponse{
			Status:  "success",
			Message: "Data fetched successfully",
		},
		Subjects:        subjects,
		LevelTrainings:  level_trainings,
		LevelExperience: level_experience,
		ClassFormats:    class_formats,
	}

	log.App.Info("Данные для поиска успешно получены")
	return response, nil
}
