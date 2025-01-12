/**
 * Страница со списком групп пользователя.
 * Отображает группы, в которых состоит пользователь.
 */
import GroupCard from "../Components/GroupCard";
import SearchBar from "../Components/SearchBar";

const MyGroups = () => {
  return (
    <div className=" w-full bg-white flex flex-col justify-start items-center pt-[15px]">
      <SearchBar
        smallMargin={true}
        placeholder="Поиск группы"
        onChange={() => {}}
      />

      <button
        className="w-[264px] bg-primary 
      text-secondary text-[20px] font-roboto font-light rounded-[12px] py-[6px] mb-[19px]"
      >
        + Создать группу
      </button>

      <div className="flex flex-col justify-start items-center gap-[18px] w-[270px]">
        <GroupCard
          name="Группа 1"
          text="Описание группы 1"
          textButton="Вступить"
          consist={true}
        />
        <GroupCard
          name="Группа 2"
          text="Описание группы 2"
          textButton="Редактировать"
        />
        <GroupCard
          name="Группа 3"
          text="Описание группы 3"
          textButton="Выйти"
        />
      </div>
    </div>
  );
};

export default MyGroups;
