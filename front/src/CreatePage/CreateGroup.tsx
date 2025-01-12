import { useState, useEffect } from "react";
import GroupDetailsCard from "./GroupDetailsCard";
import AlternationCard from "./AlternationCard";
import ScheduleRecord from "./ScheduleRecord";

/**
 * Страница создания новой группы.
 * Предоставляет форму для создания учебной группы.
 */
const CreateGroup = () => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isChecked, setIsChecked] = useState(false);
  const [date, setDate] = useState("");
  const [activeWeek, setActiveWeek] = useState("четная неделя");
  const [activeButtonWidth, setActiveButtonWidth] = useState(0);
  const [activeButtonLeft, setActiveButtonLeft] = useState(0);
  const [records, setRecords] = useState([{ id: 1, index: 1 }]);

  const handleWeekClick = (e: React.MouseEvent<HTMLButtonElement>, week: string) => {
    const button = e.currentTarget;
    setActiveButtonWidth(button.offsetWidth);
    setActiveButtonLeft(button.offsetLeft);
    setActiveWeek(week);
  };

  const addRecord = () => {
    const newIndex = records.length + 1;
    const newRecord = { id: Date.now(), index: newIndex };
    setRecords(prev => [...prev, newRecord]);
    
    // Прокрутка к кнопке после добавления записи
    setTimeout(() => {
      window.scrollTo({
        top: document.documentElement.scrollHeight,
        behavior: 'smooth'
      });
    }, 100);
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

      {/* Расписание */}
      <h2 className="text-primary text-[24px] font-light font-roboto mb-[13px] mt-[30px]">
        Расписание
      </h2>

      <div className="flex flex-col justify-start items-start bg-primary w-full mb-[150px]">
        {/* Выбор недели */}
        <div className="flex flex-row justify-center items-center gap-[20px] mt-[24px] mb-[17px] 
        w-full text-[14px] font-roboto font-light text-secondary relative">
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
          <button
            className="week-button"
            onClick={(e) => handleWeekClick(e, "все расписание")}
          >
            все расписание
          </button>
          <div
            className="absolute top-[23px] h-[2px] bg-secondary shadow-activeMenu transition-all duration-300 ease-in-out"
            style={{
              left: `${activeButtonLeft}px`,
              width: `${activeButtonWidth}px`,
            }}
          />
        </div>

        {/* Дни недели */}
        <div className="flex flex-col justify-center items-center h-[42px] pb-[15px] w-full
         text-secondary text-[20px] font-roboto font-light">
          Понедельник
        </div>

        {/* Записи */}
        <div className="w-full">
          {records.map((record, idx) => (
            <div key={record.id} 
              className="transition-all duration-300 transform origin-top"
              style={{
                animation: 'slideDown 0.3s ease-out forwards'
              }}>
              <ScheduleRecord 
                index={record.index} 
                isEven={idx % 2 === 1}
              />
            </div>
          ))}
        </div>

        {/* Кнопка добавления записи */}
        <div className="flex flex-col justify-center items-center w-full py-[8px] bg-secondary">
          <button 
            onClick={addRecord}
            className="bg-primary text-secondary text-[20px] font-roboto font-light
            py-[8px] px-[10px] rounded-[15px] hover:scale-105 transition-transform duration-300">
            + Добавить запись
          </button>
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
      `}</style>
    </div>
  );
};

export default CreateGroup;
