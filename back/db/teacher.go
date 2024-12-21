package db

import (
	"app/log"
	"app/model"
)

// GetTeachersByFilter возвращает список преподавателей по фильтру с пагинацией
func GetTeachersByFilter(filter model.TeacherFilter) ([]model.Teacher, int64, error) {
	db := App.DB

	log.App.Info("Building teacher query with filter:", filter)

	// Базовый запрос
	query := db.Model(&model.Teacher{}).
		Preload("Subjects").
		Preload("LevelTraining").
		Preload("Experience")

	if len(filter.Subjects) > 0 {
		query = query.Joins("LEFT JOIN teacher_subjects ON teachers.id = teacher_subjects.teacher_id").
			Where("teacher_subjects.subject_id IN ?", filter.Subjects).
			Group("teachers.id, teachers.created_at, teachers.updated_at, teachers.deleted_at, " +
				"teachers.user_id, teachers.name, teachers.description, teachers.education, " +
				"teachers.experience_id, teachers.price, teachers.class_format, teachers.img_url")
	}

	if len(filter.LevelTraining) > 0 {
		query = query.Joins("LEFT JOIN teacher_level_trainings ON teachers.id = teacher_level_trainings.teacher_id").
			Where("teacher_level_trainings.level_training_id IN ?", filter.LevelTraining).
			Group("teachers.id, teachers.created_at, teachers.updated_at, teachers.deleted_at, " +
				"teachers.user_id, teachers.name, teachers.description, teachers.education, " +
				"teachers.experience_id, teachers.price, teachers.class_format, teachers.img_url")
	}

	// Остальные фильтры
	if filter.Name != "" {
		query = query.Where("LOWER(teachers.name) LIKE LOWER(?)", "%"+filter.Name+"%")
	}

	if filter.Experience != 0 {
		query = query.Where("experience_id = ?", filter.Experience)
	}

	if filter.PriceFrom != nil {
		query = query.Where("price >= ?", *filter.PriceFrom)
	}
	if filter.PriceTo != nil {
		query = query.Where("price <= ?", *filter.PriceTo)
	}

	// Получаем общее количество записей
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получаем записи с пагинацией
	var teachers []model.Teacher
	if err := query.
		Order("teachers.id").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&teachers).Error; err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}
