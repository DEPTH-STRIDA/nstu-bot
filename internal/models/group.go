package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Group представляет таблицу groups
type Group struct {
	gorm.Model

	HolderID uint   `gorm:"column:holder_id"`
	Holder   User   `gorm:"foreignKey:HolderID"`
	Name     string `gorm:"column:name;not null"`
	Title    string `gorm:"column:title"`

	StartDate          *time.Time `gorm:"column:start_date"`
	IsAlternatingGroup bool       `gorm:"column:is_alternating_group;default:false"`
	EvenWeek           string     `gorm:"column:even_week;type:jsonb"`
	OddWeek            string     `gorm:"column:odd_week;type:jsonb"`
}

// TableName определяет имя таблицы для Group
func (Group) TableName() string {
	return "groups"
}

type WeekSchedule struct {
	Monday    [][]string `json:"monday,omitempty"`
	Tuesday   [][]string `json:"tuesday,omitempty"`
	Wednesday [][]string `json:"wednesday,omitempty"`
	Thursday  [][]string `json:"thursday,omitempty"`
	Friday    [][]string `json:"friday,omitempty"`
	Saturday  [][]string `json:"saturday,omitempty"`
	Sunday    [][]string `json:"sunday,omitempty"`
}

func GetScheduleOnDate(targetDate time.Time, group *Group) ([][]string, error) {
	if group == nil || group.StartDate == nil {
		return nil, errors.New("invalid group or start date")
	}

	dayOfWeek := targetDate.Weekday()

	var schedule WeekSchedule
	var weekJSON string

	if !group.IsAlternatingGroup {
		weekJSON = group.EvenWeek
	} else {
		weeks := int(targetDate.Sub(*group.StartDate).Hours() / 168) // 168 часов в неделе
		if weeks%2 == 0 {
			weekJSON = group.EvenWeek
		} else {
			weekJSON = group.OddWeek
		}
	}

	err := json.Unmarshal([]byte(weekJSON), &schedule)
	if err != nil {
		return nil, err
	}

	switch dayOfWeek {
	case time.Monday:
		return schedule.Monday, nil
	case time.Tuesday:
		return schedule.Tuesday, nil
	case time.Wednesday:
		return schedule.Wednesday, nil
	case time.Thursday:
		return schedule.Thursday, nil
	case time.Friday:
		return schedule.Friday, nil
	case time.Saturday:
		return schedule.Saturday, nil
	case time.Sunday:
		return schedule.Sunday, nil
	default:
		return nil, errors.New("invalid day of week")
	}
}

func GetDaySheduleByWeekDay(sheduleWeek string, weekDay int) ([][]string, error) {
	if sheduleWeek == "" {
		return nil, errors.New("invalid sheduleWeek")
	}

	var schedule WeekSchedule
	var weekJSON string

	err := json.Unmarshal([]byte(weekJSON), &schedule)
	if err != nil {
		return nil, err
	}

	switch weekDay {
	case 0:
		return schedule.Monday, nil
	case 1:
		return schedule.Tuesday, nil
	case 2:
		return schedule.Wednesday, nil
	case 3:
		return schedule.Thursday, nil
	case 4:
		return schedule.Friday, nil
	case 5:
		return schedule.Saturday, nil
	case 6:
		return schedule.Sunday, nil
	default:
		return nil, errors.New("invalid day of week")
	}
}

func FormatDaySchedule(schedule [][]string) string {
	str := ""

	for i, subject := range schedule {
		// Предмет должен содержать 3 поля: время, предмет, аудитория
		if len(subject) != 3 {
			continue
		}
		str += "<strong>" + fmt.Sprint(i) + ". 🕦" + subject[0] + ". 🚩" + subject[1] + "</strong> - " + subject[2]
	}

	return str
}

func FormatAllSchedule(group *Group) string {
	str := ""
	// Если не переодична
	if !group.IsAlternatingGroup {
		str += "<strong>Расписание</strong>\n\n"

	} else {
		str += "<strong>Нечетная неделя</strong>\n"
		// Обход  каждого дня недели
		for i := 0; i < 7; i++ {

			// Получение расписания на день
			shedule, err := GetDaySheduleByWeekDay(group.OddWeek, i)
			if err != nil {
				continue
			}

			// Формирование расписания
			for i, subject := range shedule {
				// Предмет должен содержать 3 поля: время, предмет, аудитория
				if len(subject) != 3 {
					continue
				}
				str += "<strong>" + fmt.Sprint(i) + ". 🕦" + subject[0] + ". 🚩" + subject[1] + "</strong> - " + subject[2]
			}
		}

	}
	if group.IsAlternatingGroup {
		str += "<strong>Четная неделя</strong>\n"
	}
	// Обход  каждого дня недели
	for i := 0; i < 7; i++ {

		// Получение расписания на день
		shedule, err := GetDaySheduleByWeekDay(group.EvenWeek, i)
		if err != nil {
			continue
		}

		// Формирование расписания
		for i, subject := range shedule {
			// Предмет должен содержать 3 поля: время, предмет, аудитория
			if len(subject) != 3 {
				continue
			}
			str += "<strong>" + fmt.Sprint(i) + ". 🕦" + subject[0] + ". 🚩" + subject[1] + "</strong> - " + subject[2]
		}
	}

	return str
}

func GetDayOfWeekInRussian(date time.Time) string {
	daysOfWeek := []string{
		"воскресенье",
		"понедельник",
		"вторник",
		"среда",
		"четверг",
		"пятница",
		"суббота",
	}

	// time.Weekday() возвращает число от 0 (воскресенье) до 6 (суббота)
	dayIndex := int(date.Weekday())

	return daysOfWeek[dayIndex]
}
