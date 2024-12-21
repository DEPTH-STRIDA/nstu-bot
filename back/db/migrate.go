package db

import (
	"app/log"
	"app/model"
	"fmt"
	"math/rand"
)

func UpdateTeacherPhotos() error {
	// Получаем всех учителей
	var teachers []model.Teacher
	if err := App.Find(&teachers).Error; err != nil {
		return fmt.Errorf("failed to fetch teachers: %v", err)
	}

	// Для каждого учителя генерируем случайную ссылку на фото
	for _, teacher := range teachers {
		randomNum := rand.Intn(4) + 1 // Генерируем число от 1 до 4
		imgUrl := fmt.Sprintf("/img/%d.jpg", randomNum)

		if err := App.Model(&teacher).Update("img_url", imgUrl).Error; err != nil {
			return fmt.Errorf("failed to update teacher %d: %v", teacher.ID, err)
		}
	}

	fmt.Printf("Updated %d teacher records\n", len(teachers))
	return nil
}

// Migrate выполняет миграции базы данных
func (db *DataBase) Migrate() error {

	log.App.Info("Starting database migration...")

	// Сначала создаем все таблицы
	log.App.Info("Running auto-migration...")
	err := db.AutoMigrate(
		&model.Subject{},
		&model.LevelTraining{},
		&model.User{},
		&model.Chat{},
		&model.Message{},
		&model.Admin{},
		&model.Student{},
		&model.Teacher{},
		&model.TutorExperience{},
		&model.Service{},
		&model.ClassFormat{},
	)
	if err != nil {
		log.App.Error("Auto-migration failed:", err)
		return err
	}
	log.App.Info("Auto-migration completed successfully")

	App.Create(&model.ClassFormat{Name: "Индивидуальные"})
	App.Create(&model.ClassFormat{Name: "Групповые"})

	// Создаем базовые данные
	subjects := []model.Subject{
		{Name: "Математика"},
		{Name: "Физика"},
		{Name: "Химия"},
		{Name: "Биология"},
		{Name: "История"},
		{Name: "Обществознание"},
		{Name: "Русский язык"},
		{Name: "Литература"},
		{Name: "Английский язык"},
		{Name: "Информатика"},
	}

	levels := []model.LevelTraining{
		{Name: "Начальное образование"},
		{Name: "5-6 класс"},
		{Name: "7-8 класс"},
		{Name: "9 класс"},
		{Name: "10-11 класс"},
		{Name: "Подготовка к ОГЭ"},
		{Name: "Подготовка к ЕГЭ"},
	}

	experiences := []model.TutorExperience{
		{Name: "Без опыта", Value: 0},
		{Name: "До 1 года", Value: 1},
		{Name: "1-3 года", Value: 2},
		{Name: "3-5 лет", Value: 3},
		{Name: "5-10 лет", Value: 4},
		{Name: "Более 10 лет", Value: 5},
	}

	services := []model.Service{
		{Name: "Подготовка к ЕГЭ"},
		{Name: "Подготовка к ОГЭ"},
		{Name: "Помощь с домашними заданиями"},
		{Name: "Подготовка к олимпиадам"},
		{Name: "Подготовка к контрольным"},
		{Name: "Устранение пробелов в знаниях"},
		{Name: "Развитие интереса к предмету"},
		{Name: "Углубленное изучение предмета"},
	}

	// Создаем записи
	log.App.Info("Creating subjects...")
	for _, s := range subjects {
		var exists bool
		db.Model(&model.Subject{}).Where("name = ?", s.Name).Select("count(*) > 0").Find(&exists)
		if !exists {
			log.App.Info(fmt.Sprintf("Creating subject: %s", s.Name))
			if err := db.Create(&s).Error; err != nil {
				log.App.Warn(fmt.Sprintf("Failed to create subject %s: %v", s.Name, err))
				continue
			}
		} else {
			log.App.Info(fmt.Sprintf("Subject already exists: %s", s.Name))
		}
	}

	log.App.Info("Creating levels...")
	for _, l := range levels {
		var exists bool
		db.Model(&model.LevelTraining{}).Where("name = ?", l.Name).Select("count(*) > 0").Find(&exists)
		if !exists {
			log.App.Info(fmt.Sprintf("Creating level: %s", l.Name))
			if err := db.Create(&l).Error; err != nil {
				log.App.Warn(fmt.Sprintf("Failed to create level %s: %v", l.Name, err))
				continue
			}
		}
	}

	log.App.Info("Creating experiences...")
	for _, e := range experiences {
		var exists bool
		db.Model(&model.TutorExperience{}).Where("name = ?", e.Name).Select("count(*) > 0").Find(&exists)
		if !exists {
			log.App.Info(fmt.Sprintf("Creating experience: %s", e.Name))
			if err := db.Create(&e).Error; err != nil {
				log.App.Warn(fmt.Sprintf("Failed to create experience %s: %v", e.Name, err))
				continue
			}
		}
	}

	log.App.Info("Creating services...")
	for _, s := range services {
		var exists bool
		db.Model(&model.Service{}).Where("name = ?", s.Name).Select("count(*) > 0").Find(&exists)
		if !exists {
			log.App.Info(fmt.Sprintf("Creating service: %s", s.Name))
			if err := db.Create(&s).Error; err != nil {
				log.App.Warn(fmt.Sprintf("Failed to create service %s: %v", s.Name, err))
				continue
			}
		}
	}

	log.App.Info("Base data created successfully")

	// Создаем преподавателей
	log.App.Info("Creating teachers...")
	for i := 0; i < 50; i++ {
		log.App.Info(fmt.Sprintf("Creating teacher %d/50", i+1))
		user := model.User{
			Name:     fmt.Sprintf("Teacher%d Name%d", i, rand.Intn(1000)),
			Email:    fmt.Sprintf("teacher%d@test.com", i),
			Password: "$2a$10$Yw3hJ6u0RM8SKx1bBxk9Q.RZpKKTPh3f8Q9E6qZjY9hMAqv0Jy7C2",
			Role:     model.TeacherRole,
		}

		if err := db.Create(&user).Error; err != nil {
			log.App.Error("Failed to create teacher user:", err)
			return err
		}
		log.App.Info(fmt.Sprintf("Created teacher user: %s", user.Name))

		// Выбираем случайные предметы и уровни
		var selectedSubjects []model.Subject
		if err := db.Order("RANDOM()").Limit(rand.Intn(3) + 1).Find(&selectedSubjects).Error; err != nil {
			return err
		}

		var selectedLevels []model.LevelTraining
		if err := db.Order("RANDOM()").Limit(rand.Intn(3) + 1).Find(&selectedLevels).Error; err != nil {
			return err
		}

		// Получаем все форматы
		var formats []model.ClassFormat
		if err := db.Find(&formats).Error; err != nil {
			return err
		}

		teacher := model.Teacher{
			UserID:        user.ID,
			Name:          user.Name,
			Description:   fmt.Sprintf("Description for teacher %d", i),
			Education:     fmt.Sprintf("University %d", rand.Intn(100)),
			ExperienceID:  uint(rand.Intn(5) + 1),
			Price:         rand.Intn(3000) + 1000,
			ImgUrl:        "/img/EygRcOLGR1s.jpg",
			Subjects:      selectedSubjects,
			LevelTraining: selectedLevels,
			// Случайный выбор форматов
			ClassFormats: formats[:rand.Intn(len(formats)+1)], // случайное количество форматов (0 до 2)
		}

		// При создании преподавателей добавляем случайные услуги
		var selectedServices []model.Service
		if err := db.Order("RANDOM()").Limit(rand.Intn(3) + 1).Find(&selectedServices).Error; err != nil {
			return err
		}
		teacher.Services = selectedServices

		if err := db.Create(&teacher).Error; err != nil {
			return err
		}

		if err := db.Model(&user).Update("TeacherID", teacher.ID).Error; err != nil {
			return err
		}

		log.App.Info(fmt.Sprintf("Created teacher profile for: %s", user.Name))
	}

	// Создаем админов
	log.App.Info("Creating admins...")
	for i := 0; i < 5; i++ {
		log.App.Info(fmt.Sprintf("Creating admin %d/5", i+1))
		user := model.User{
			Name:     fmt.Sprintf("Admin%d", i),
			Email:    fmt.Sprintf("admin%d@test.com", i),
			Password: "$2a$10$Yw3hJ6u0RM8SKx1bBxk9Q.RZpKKTPh3f8Q9E6qZjY9hMAqv0Jy7C2",
			Role:     model.AdminRole,
		}

		if err := db.Create(&user).Error; err != nil {
			log.App.Error("Failed to create admin user:", err)
			return err
		}
		log.App.Info(fmt.Sprintf("Created admin user: %s", user.Name))

		admin := model.Admin{
			Description: fmt.Sprintf("Description for admin %d", i),
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}

		if err := db.Model(&user).Update("AdminID", admin.ID).Error; err != nil {
			return err
		}

		log.App.Info(fmt.Sprintf("Created admin profile for: %s", user.Name))
	}

	// Создаем учеников
	log.App.Info("Creating students...")
	for i := 0; i < 5; i++ {
		log.App.Info(fmt.Sprintf("Creating student %d/5", i+1))
		user := model.User{
			Name:     fmt.Sprintf("Student%d", i),
			Email:    fmt.Sprintf("student%d@test.com", i),
			Password: "$2a$10$Yw3hJ6u0RM8SKx1bBxk9Q.RZpKKTPh3f8Q9E6qZjY9hMAqv0Jy7C2",
			Role:     model.StudentRole,
		}

		if err := db.Create(&user).Error; err != nil {
			log.App.Error("Failed to create student user:", err)
			return err
		}
		log.App.Info(fmt.Sprintf("Created student user: %s", user.Name))

		// Выбираем случайные предметы для ученика
		var selectedSubjects []model.Subject
		if err := db.Order("RANDOM()").Limit(rand.Intn(3) + 1).Find(&selectedSubjects).Error; err != nil {
			return err
		}

		student := model.Student{
			Subjects:    selectedSubjects,
			Class:       fmt.Sprintf("%d класс", rand.Intn(11)+1),
			Description: fmt.Sprintf("Description for student %d", i),
		}

		if err := db.Create(&student).Error; err != nil {
			return err
		}

		if err := db.Model(&user).Update("StudentID", student.ID).Error; err != nil {
			return err
		}

		log.App.Info(fmt.Sprintf("Created student profile for: %s", user.Name))
	}

	log.App.Info("Migration completed successfully")
	UpdateTeacherPhotos()
	return nil
}
