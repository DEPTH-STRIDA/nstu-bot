import { useState, useEffect } from "react";
import GroupDetailsCard from "./GroupDetailsCard";
import AlternationCard from "./AlternationCard";
import ScheduleRecord from "./ScheduleRecord";

const DAYS = {
  "Понедельник": "Понедельник",
  "Вторник": "Вторник",
  "Среда": "Среда",
  "Четверг": "Четверг",
  "Пятница": "Пятница",
  "Суббота": "Суббота"
} as const;

type DayKey = keyof typeof DAYS;
type Schedule = {[key in DayKey]: { id: number; index: number }[]};

const CreateGroup = () => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isChecked, setIsChecked] = useState(false);
  const [date, setDate] = useState("");
  const [activeWeek, setActiveWeek] = useState("четная неделя");
  const [activeButtonWidth, setActiveButtonWidth] = useState(0);
  const [activeButtonLeft, setActiveButtonLeft] = useState(0);
  
  // Отдельные состояния для четной и нечетной недели
  const [evenSchedule, setEvenSchedule] = useState<Schedule>(() => 
    Object.keys(DAYS).reduce((acc, day) => ({
      ...acc,
      [day]: []
    }), {} as Schedule)
  );
  
  const [oddSchedule, setOddSchedule] = useState<Schedule>(() => 
    Object.keys(DAYS).reduce((acc, day) => ({
      ...acc,
      [day]: []
    }), {} as Schedule)
  );

  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [movingId, setMovingId] = useState<number | null>(null);

  // Получаем текущее расписание на основе активной недели
  const getCurrentSchedule = () => {
    if (!isChecked) return evenSchedule;
    return activeWeek === "четная неделя" ? evenSchedule : oddSchedule;
  };

  const setCurrentSchedule = (newSchedule: Schedule) => {
    if (!isChecked) {
      setEvenSchedule(newSchedule);
      return;
    }
    if (activeWeek === "четная неделя") {
      setEvenSchedule(newSchedule);
    } else {
      setOddSchedule(newSchedule);
    }
  };

  // При изменении состояния галочки
  useEffect(() => {
    if (!isChecked) {
      setActiveWeek("расписание");
    } else if (activeWeek === "расписание") {
      setActiveWeek("четная неделя");
    }
  }, [isChecked]);

  const handleWeekClick = (e: React.MouseEvent<HTMLButtonElement>, week: string) => {
    const button = e.currentTarget;
    setActiveButtonWidth(button.offsetWidth);
    setActiveButtonLeft(button.offsetLeft);
    setActiveWeek(week);
  };

  const addRecord = (day: DayKey) => {
    setCurrentSchedule({
      ...getCurrentSchedule(),
      [day]: [...getCurrentSchedule()[day], { id: Date.now(), index: getCurrentSchedule()[day].length + 1 }]
    });
    
    setTimeout(() => {
      window.scrollTo({
        top: document.documentElement.scrollHeight,
        behavior: 'smooth'
      });
    }, 100);
  };

  const deleteRecord = (day: DayKey, id: number) => {
    setDeletingId(id);
    setTimeout(() => {
      setCurrentSchedule({
        ...getCurrentSchedule(),
        [day]: getCurrentSchedule()[day]
          .filter(record => record.id !== id)
          .map((record, idx) => ({ ...record, index: idx + 1 }))
      });
      setDeletingId(null);
    }, 300);
  };

  const moveRecord = (day: DayKey, id: number, direction: 'up' | 'down') => {
    const dayRecords = getCurrentSchedule()[day];
    const currentIndex = dayRecords.findIndex(record => record.id === id);
    if (
      (direction === 'up' && currentIndex === 0) || 
      (direction === 'down' && currentIndex === dayRecords.length - 1)
    ) return;

    setMovingId(id);
    setTimeout(() => {
      const newSchedule = { ...getCurrentSchedule() };
      const newDayRecords = [...newSchedule[day]];
      const swapIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
      
      [newDayRecords[currentIndex], newDayRecords[swapIndex]] = 
      [newDayRecords[swapIndex], newDayRecords[currentIndex]];
      
      newSchedule[day] = newDayRecords.map((record, idx) => ({
        ...record,
        index: idx + 1
      }));

      setCurrentSchedule(newSchedule);
      setMovingId(null);
    }, 300);
  };

  useEffect(() => {
    const buttons = document.querySelectorAll(".week-button");
    const activeButton = Array.from(buttons).find(
      (button) => button.textContent === activeWeek
    ) as HTMLButtonElement | undefined;

    if (activeButton) {
      setActiveButtonWidth(activeButton.offsetWidth);
      setActiveButtonLeft(activeButton.offsetLeft);
    }
  }, [activeWeek]);

  const currentSchedule = getCurrentSchedule();

  return (
    <div className="w-full bg-white flex flex-col justify-start items-center pt-[23px]">
      {/* Ввод названия и описания группы */}
      <GroupDetailsCard
        name={name}
        description={description}
        onNameChange={setName}
        onDescriptionChange={setDescription}
      />

      {/* Чредование */}
      <AlternationCard
        onRetrieveCheckboxState={setIsChecked}
        onRetrieveDate={setDate}
      />

      <div className="flex flex-col justify-start items-start bg-primary w-full mt-[22px]">
        <div className="flex flex-row justify-center items-center gap-[20px] mt-[13px] mb-[13px] 
        w-full text-[14px] font-roboto font-light text-secondary relative">
          {isChecked ? (
            <>
              <button
                className="week-button"
                onClick={(e) => handleWeekClick(e, "четная неделя")}
              >
                четная неделя
              </button>
              <button
                className="week-button"
                onClick={(e) => handleWeekClick(e, "нечетная неделя")}
              >
                нечетная неделя
              </button>
            </>
          ) : (
            <button
              className="week-button"
              onClick={(e) => handleWeekClick(e, "расписание")}
            >
              расписание
            </button>
          )}
          <div
            className="absolute top-[23px] h-[2px] bg-secondary shadow-activeMenu transition-all duration-300 ease-in-out"
            style={{
              left: `${activeButtonLeft}px`,
              width: `${activeButtonWidth}px`,
            }}
          />
        </div>

        {/* Дни недели с расписаниями */}
        <div className="w-full transition-opacity duration-300">
          {(Object.keys(DAYS) as DayKey[]).map((day) => (
            <div key={day} className="w-full">
              <div className="flex flex-col justify-center items-center h-[42px] pb-[15px] w-full
                text-secondary text-[20px] font-roboto font-light">
                {DAYS[day]}
              </div>

              {/* Записи для дня */}
              <div className="w-full">
                {currentSchedule[day].map((record, idx) => (
                  <div key={record.id} 
                    className={`transition-all duration-300 transform origin-top
                      ${deletingId === record.id ? 'opacity-0 scale-y-0' : 'opacity-100 scale-y-100'}
                      ${movingId === record.id ? 'z-10 shadow-lg' : ''}`}
                    style={{
                      animation: deletingId === record.id ? 'slideUp 0.3s ease-out forwards' :
                                movingId === record.id ? 'none' : 'slideDown 0.3s ease-out forwards'
                    }}>
                    <ScheduleRecord 
                      index={record.index} 
                      isEven={idx % 2 === 1}
                      onDelete={() => deleteRecord(day, record.id)}
                      onMoveUp={() => moveRecord(day, record.id, 'up')}
                      onMoveDown={() => moveRecord(day, record.id, 'down')}
                      canMoveUp={idx !== 0}
                      canMoveDown={idx !== currentSchedule[day].length - 1}
                    />
                  </div>
                ))}
              </div>

              {/* Кнопка добавления записи для дня */}
              <div className="flex flex-col justify-center items-center w-full py-[8px] bg-secondary">
                <button 
                  onClick={() => addRecord(day)}
                  className="bg-primary text-secondary text-[20px] font-roboto font-light
                  py-[8px] px-[10px] rounded-[15px] hover:scale-105 transition-transform duration-300">
                  + Добавить запись
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <style>{`
        @keyframes slideDown {
          from {
            opacity: 0;
            transform: translateY(-10px) scaleY(0.5);
          }
          to {
            opacity: 1;
            transform: translateY(0) scaleY(1);
          }
        }
        @keyframes slideUp {
          from {
            opacity: 1;
            transform: translateY(0) scaleY(1);
          }
          to {
            opacity: 0;
            transform: translateY(-10px) scaleY(0);
          }
        }
      `}</style>
    </div>
  );
};

export default CreateGroup;
