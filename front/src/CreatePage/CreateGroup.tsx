import { useState } from "react";
import GroupDetailsCard from "./GroupDetailsCard";
import AlternationCard from "./AlternationCard";

/**
 * Страница создания новой группы.
 * Предоставляет форму для создания учебной группы.
 */
const CreateGroup = () => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isChecked, setIsChecked] = useState(false);
  const [date, setDate] = useState("");

  return (
    <div className=" w-full bg-white flex flex-col justify-start items-center pt-[23px] gap-[23px]">
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
      <h2 className="text-primary text-[24px] font-light font-roboto">
        Расписание
      </h2>

      <div className="flex flex-col justify-start items-start bg-primary">
        <div className="flex flex-row justify-start items-start gap-[10px] mt-[24px] mb-[17px]
         text-secondary">
          <button>четная неделя</button>
          <button>нечетная неделя</button>
          <button>все расписание</button>
        </div>
      </div>
    </div>
  );
};

export default CreateGroup;
