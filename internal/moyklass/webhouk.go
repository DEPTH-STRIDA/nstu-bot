package moyklass

import (
	"encoding/json"
	"time"
)

type ObjectInterface interface {
	New() ObjectInterface
}

var EventTypes = map[string]ObjectInterface{
	// Обьект учеников
	"user_new":           User{}, // Ученик создан
	"user_changed":       User{}, // Информация об ученике изменилась
	"user_changed_state": User{}, // Статус ученика изменился
	"user_birthday":      User{}, // Ученик отмечает день рождения

	// Объект записи в группу
	"join_new":           Join{}, // Создана запись в группу
	"join_changed":       Join{}, // Информация о записи в группу изменилась
	"join_changed_state": Join{}, // Статус записи в группу изменился

	// Объект удаленной записи в группу:
	"join_deleted": Exit{}, // Удалена запись в группу

	// Объект группы
	"class_new":         Group{}, // Группа создана
	"class_changed":     Group{}, // Информация о группе изменилась
	"class_deleted":     Group{}, // Группа удалена
	"class_start_days":  Group{}, // Старт группы менее чем через N дней
	"class_start_hours": Group{}, // Старт группы менее чем через N часов

	// Объект занятия
	"lesson_changed":     Lesson{}, // Информация о занятии изменилась
	"lesson_deleted":     Lesson{}, // Занятие удалено
	"lesson_start_days":  Lesson{}, // Занятие начинается менее чем через N дней
	"lesson_start_hours": Lesson{}, // Занятие начинается через N часов

	// Объект платежа
	"payment_new": Payment{}, // Принят платеж

	// Объект списания
	"debit_new": Debit{}, // Новое списание у ученика

	// Объект записи на занятие
	"sub_lesson_in_debt":    Subscription{}, // Занятие для записи проведено “в долг”
	"lesson_record_new":     Subscription{}, // Создана запись на занятие
	"lesson_record_changed": Subscription{}, //Информация о записи на занятие изменилась

	// Объект отмененной записи
	"lesson_record_deleted": CancelSubscription{}, // Запись на занятие удалена

	// Объект абонемента пользователя
	"sub_days_next_payment": SeasonTicket{}, // Осталось менее N дней до внесения очередного платежа/долга по абонементу
	"sub_lessons_left":      SeasonTicket{}, // В абонементе осталось менее N посещений
	"sub_end_days":          SeasonTicket{}, // Абонемент заканчивается менее чем через N дней
	"sub_new":               SeasonTicket{}, // Создание нового абонемента

	// Объект двух пропущенных занятий подряд
	"user_consecutive_visit_missed_2": TwoMissedClassesRow{}, // Ученик имеет два пропуска подряд без уважительной причины

	// События оценок
	"lesson_mark_set":     Mark{}, // Проставлена оценка
	"lesson_mark_deleted": Mark{}, // Оценка удалена

	// События заданий на уроке
	"lesson_task_new":     TaskHomeWork{}, // Задание создано
	"lesson_task_changed": TaskHomeWork{}, // Задание изменено

	// События задач
	"task_new":     Task{},        // Задача создана
	"task_changed": Task{},        // Задача изменена
	"task_deleted": DeletedTask{}, // Задача удалена
}

type BaseUpdate struct {
	CompanyId     int64         `json:"companyId"`     // Id компании
	Event         string        `json:"event"`         // Код события
	Time          time.Time     `json:"time"`          // Отметка времени UNIX
	ChangedFields ChangedFields `json:"changedFields"` // Массив измененных полей
	Init          Init          `json:"init"`          // Объект инициатора события
}

type Update struct {
	Object json.RawMessage `json:"object"`
	BaseUpdate
}

// ChangedFields содержит массив измененных полей.
//
//	user_changed - Информация об ученике изменилась
//	user_changed_state - Статус ученика изменился
//	join_changed - Информация о записи в группу изменилась
//	join_changed_state - Статус записи в группу изменился
//	class_changed - Информация о группе изменилась
//	lesson_changed - Информация о занятии изменилась
//	lesson_record_changed - Информация о записи на занятие изменилась
//	lesson_task_changed - Информация о задании изменилась
type ChangedFields []string

// Init содержит информацию о том, кто инициировал событие.
// Поле from может принимать значения:
//
//	CRM - событие вызвано из CRM, в init добавляется managerId - ID менеджера
//	import - событие вызвано из импорта, в init добавляется managerId - ID менеджера
//	LK - событие вызвано из Личного Кабинета ученика, в init добавляется userId - ID ученика
//	widget - событие вызвано из Виджетов, в init добавляется widgetId - ID виджета
//	integration - событие вызвано из Интеграции, в init добавляется integrationId - ID интеграции
//	API - событие вызвано из API
//	SYSTEM - событие сгенерировано автоматически.
type Init struct {
	From      string `json:"from"`      // CRM, import, LK, widget, integration, API, SYSTEM
	ManagerID int64  `json:"managerId"` // ID менеджера (CRM, import)
}

// User обьект учеников. используется в событиях: user_new, user_changed, user_changed_state, user_birthday
type User struct {
	UserID             int           `json:"userId"`             // ID ученика
	Filials            []int         `json:"filials"`            // ID филиалов ученика
	Attributes         []interface{} `json:"attributes"`         // Дополнительные атрибуты ученика
	Name               string        `json:"name"`               // Полное имя ученика
	Balans             int           `json:"balans"`             // Баланс ученика
	Email              string        `json:"email"`              // Email ученика
	CreatedAt          string        `json:"createdAt"`          // Дата создания ученика
	StatusID           int           `json:"statusId"`           // ID статуса ученика
	AdvSourceID        int           `json:"advSourceId"`        // ID информационного источника (откуда ученик узнал о компании)
	ResponsibleIDs     []int         `json:"responsibleIds"`     // Список ID ответственных менеджеров
	StatusReasonID     int           `json:"statusReasonId"`     // ID причины смены статуса ученика
	CreateSourceID     int           `json:"createSourceId"`     // ID источника создания ученика
	PrevStatusID       int           `json:"prevStatusId"`       // ID предыдущего статуса ученика
	PrevStatusReasonID int           `json:"prevStatusReasonId"` // ID причины смены предыдущего статуса ученика
	Phone              string        `json:"phone"`              // Номер телефона ученика
	Remind             string        `json:"remind"`             // Напоминание
}

func (User) New() ObjectInterface {
	return User{}
}

// Join объект записи в группу. используется в событиях: join_new, join_changed, join_changed_state
type Join struct {
	JoinID                   int         `json:"joinId"`                   // ID заявки
	Price                    int         `json:"price"`                    // Цена (для групп с разовой оплатой)
	UserID                   int         `json:"userId"`                   // ID ученика
	ClassID                  int         `json:"classId"`                  // ID группы
	Comment                  string      `json:"comment"`                  // Комментарий
	AutoJoin                 bool        `json:"autoJoin"`                 // Автоматически записывать в статусе "Учится" на все занятия в группе
	ManagerID                int         `json:"managerId"`                // ID ответственного сотрудника
	RemindSum                int         `json:"remindSum"`                // Сумма долга к оплате
	CreatedAt                string      `json:"createdAt"`                // Дата создания
	StatusID                 int         `json:"statusId"`                 // ID статуса заявки
	AdvSourceID              int         `json:"advSourceId"`              // ID информационного источника (откуда ученик узнал о компании)
	PrevStatusID             int         `json:"prevStatusId"`             // ID предыдущего статуса заявки
	CreateSourceID           int         `json:"createSourceId"`           // ID источника создания
	StatusChangeReasonID     int         `json:"statusChangeReasonId"`     // ID причины смены статуса
	PrevStatusChangeReasonID int         `json:"prevStatusChangeReasonId"` // ID причины смены предыдущего статуса
	RemindDate               string      `json:"remindDate"`               // Срок оплаты долга
	Stats                    interface{} `json:"stats"`                    // Статистика по записи
}

func (Join) New() ObjectInterface {
	return Join{}
}

// Exit объект удаленной записи в группу. Используется в событии join_deleted
type Exit struct {
	UserId  int `json:"userId"`  // ID ученика
	ClassId int `json:"classId"` // ID группы
	JoinId  int `json:"joinId"`  // ID заявки
}

func (Exit) New() ObjectInterface {
	return Exit{}
}

// Group объект группы. Используется в событиях: class_new, class_changed, class_deleted, class_start_days, class_start_hours
type Group struct {
	ClassID        int    `json:"classId"`        // ID группы
	Name           string `json:"name"`           // Название
	TeacherIDs     []int  `json:"teacherIds"`     // Список ID ведущих группы
	Color          string `json:"color"`          // Цвет
	Price          int    `json:"price"`          // Цена
	Comment        string `json:"comment"`        // Комментарий
	StatusID       int    `json:"statusId"`       // ID статуса
	CourseID       int    `json:"courseId"`       // ID программы
	CourseName     string `json:"courseName"`     // название программы
	ShowDates      bool   `json:"showDates"`      // Отображение даты начала у названия группы
	MaxStudents    int    `json:"maxStudents"`    // Максимальное количество студентов
	FilialID       int    `json:"filialId"`       // ID филиала
	PriceComment   string `json:"priceComment"`   // Комментарий к цене
	PriceForWidget string `json:"priceForWidget"` // Цена для виджетов
	CreatedAt      string `json:"createdAt"`      // Дата создания
	PayType        string `json:"payType"`        // Способ оплаты за обучение. full - разово, lessons - за занятия
	BeginDate      string `json:"beginDate"`      // Старт занятий
}

func (Group) New() ObjectInterface {
	return Group{}
}

// Lesson объект занятия. Используется в событиях: lesson_changed, lesson_deleted, lesson_start_days, lesson_start_hours
type Lesson struct {
	LessonID    int          `json:"lessonId"`    // ID занятия
	Topic       string       `json:"topic"`       // Тема занятия
	EndTime     string       `json:"endTime"`     // Время окончания занятия
	ClassID     int          `json:"classId"`     // ID группы
	Status      int          `json:"status"`      // Статус
	CreatedAt   string       `json:"createdAt"`   // Дата создания
	BeginTime   string       `json:"beginTime"`   // Время начала
	Description string       `json:"description"` // Описание занятия
	FilialID    int          `json:"filialId"`    // ID филиала
	RoomID      int          `json:"roomId"`      // ID аудитории
	Date        string       `json:"date"`        // Дата начала занятия
	Paramss     LessonParams `json:"params"`      // Параметры занятия
}

// LessonParams содержит дополнительные параметры занятия
type LessonParams struct {
	Webinars []interface{} `json:"webinars"` // Массив объектов вебинаров
	Videos   []interface{} `json:"videos"`   // Массив объектов видео
}

func (Lesson) New() ObjectInterface {
	return Lesson{}
}

// Payment объект платежа. Используется в событиях: payment_new
type Payment struct {
	UserID             int    `json:"userId"`             // ID ученика
	PaymentID          int    `json:"paymentId"`          // ID платежа
	Summa              int    `json:"summa"`              // Сумма платежа
	UserSubscriptionID int    `json:"userSubscriptionId"` // ID абонемента ученика
	PaymentTypeID      int    `json:"paymentTypeId"`      // ID типа оплаты
	Date               string `json:"date"`               // Дата оплаты
}

func (Payment) New() ObjectInterface {
	return Payment{}
}

// Debit объект списания. Используется в событиях: debit_new
type Debit struct {
	UserID             int    `json:"userId"`             // ID ученика
	PaymentID          int    `json:"paymentId"`          // ID платежа
	InvoiceID          int    `json:"invoiceId"`          // ID счета
	Summa              int    `json:"summa"`              // Сумма списания
	ClassID            int    `json:"classId"`            // ID группы
	SaleID             int    `json:"saleId"`             // ID продажи товара
	UserSubscriptionID int    `json:"userSubscriptionId"` // ID абонемента ученика
	Date               string `json:"date"`               // Дата списания
}

func (Debit) New() ObjectInterface {
	return Debit{}
}

// Subscription Объект записи на занятие. Используется в событиях: sub_lesson_in_debt, lesson_record_new, lesson_record_changed.
type Subscription struct {
	LessonRecordID int       `json:"lessonRecordId"` // ID записи на занятие
	Free           bool      `json:"free"`           // Бесплатная запись
	Test           bool      `json:"test"`           // Пробная запись на занятие
	Skip           bool      `json:"skip"`           // Не учитывать запись в количестве занятых мест
	Visit          bool      `json:"visit"`          // Статус посещения
	UserID         int       `json:"userId"`         // ID ученика
	LessonID       int       `json:"lessonId"`       // ID занятия
	CreatedAt      string    `json:"createdAt"`      // Дата создания
	GoodReason     bool      `json:"goodReason"`     // Уважительная причина отсутствия
	JoinStats      JoinStats `json:"joinStats"`      // Статистика записи в группу
}

// JoinStats содержит статистику посещений
type JoinStats struct {
	LessonVisited bool `json:"lesson_visited"` // Ученик посетил занятие
	LessonMissed  bool `json:"lesson_missed"`  // Ученик пропустил занятие
}

func (Subscription) New() ObjectInterface {
	return Subscription{}
}

// CanselSubscription объект отмененной записи. Используется в событиях: lesson_record_deleted.
type CancelSubscription struct {
	classId  int // ID группы
	userId   int // ID ученика
	lessonId int //ID занятия
}

func (CancelSubscription) New() ObjectInterface {
	return CancelSubscription{}
}

// SeasonTicket объект абонемента пользователя. Используется в событиях: sub_days_next_payment, sub_lessons_left, sub_end_days, sub_new.
type SeasonTicket struct {
	UserSubscriptionID int    `json:"userSubscriptionId"` // ID абонемента
	UserID             int    `json:"userId"`             // ID ученика
	StatusID           int    `json:"statusId"`           // ID статуса абонемента
	VisitedCount       int    `json:"visitedCount"`       // Количество списанных занятий в абонементе
	VisitCount         int    `json:"visitCount"`         // Количество занятий в абонементе
	RemindSum          int    `json:"remindSumm"`         // Сумма долга к оплате
	EndDate            string `json:"endDate"`            // Дата окончания действия
	SellDate           string `json:"sellDate"`           // Дата продажи
	BeginDate          string `json:"beginDate"`          // Дата начала действия
	RemindDate         string `json:"remindDate"`         // Срок оплаты долга
}

func (SeasonTicket) New() ObjectInterface {
	return SeasonTicket{}
}

// TwoMissedClassesRow объект двух пропущенных занятий подряд. Используется в событиях: user_consecutive_visit_missed_2
type TwoMissedClassesRow struct {
	ClassID   int   `json:"classId"`   // ID группы
	UserID    int   `json:"userId"`    // ID ученика
	LessonIDs []int `json:"lessonIds"` // ID последних двух пропущенных занятий
}

func (TwoMissedClassesRow) New() ObjectInterface {
	return TwoMissedClassesRow{}
}

// Mark Объект оценки. Используется в событиях: lesson_mark_set, lesson_mark_deleted
type Mark struct {
	LessonID int    `json:"lessonId"` // ID занятия
	UserID   int    `json:"userId"`   // ID ученика
	Type     string `json:"type"`     // Тип оценки, home либо lesson
	Value    int    `json:"value"`    // Оценка
}

func (Mark) New() ObjectInterface {
	return Mark{}
}

// TaskHomeWork объект задания. Используется в событиях: lesson_task_new, lesson_task_changed
type TaskHomeWork struct {
	LessonID   int    `json:"lessonId"`   // ID занятия
	Files      []int  `json:"files"`      // Список ID прикрепленных файлов
	Type       string `json:"type"`       // Тип задания, home либо lesson
	ShowToUser bool   `json:"showToUser"` // Доступно ли задание ученику для просмотра
}

func (TaskHomeWork) New() ObjectInterface {
	return TaskHomeWork{}
}

// Task Объект задачи. Используется в событиях: task_new - Задача создана, task_changed - задача изменена
type Task struct {
	LessonID   int    `json:"lessonId"`   // ID задачи
	Body       string `json:"body"`       // Текст задачи
	CompanyID  int    `json:"companyId"`  // ID компании
	BeginDate  string `json:"beginDate"`  // Начало задачи
	EndDate    string `json:"endDate"`    // Окончание задачи
	IsAllDay   bool   `json:"isAllDay"`   // Задача на весь день
	IsComplete bool   `json:"isComplete"` // Задача выполнена
	UserID     int    `json:"userId"`     // ID ученика
	ClassIDs   []int  `json:"classIds"`   // Список связанных групп
	FilialIDs  []int  `json:"filialIds"`  // Список связанных филиалов
	CategoryID int    `json:"categoryId"` // ID категории задачи
	InvoiceID  int    `json:"invoiceId"`  // ID счета
	ManagerIDs []int  `json:"managerIds"` // Список id ответственных сотрудников
}

func (Task) New() ObjectInterface {
	return Task{}
}

// DeletedTask объект удаленной задачи. Используется в событиях: task_deleted - Задание удалено
type DeletedTask struct {
	TaskID int `json:"taskId"` // ID задачи
}

func (DeletedTask) New() ObjectInterface {
	return DeletedTask{}
}
