/**
 * Страница со списком групп пользователя.
 * Отображает группы, в которых состоит пользователь.
 */
import GroupCard from "../Components/GroupCard";

const MyGroups = () => {
  return (
    <div className=" w-full  flex flex-col justify-center items-center gap-[18px] py-[30px]">
      <GroupCard name="Группа 1" text="Описание группы 1" textButton="Вступить" consist={true} />
      <GroupCard name="Группа 2" text="Описание группы 2" textButton="Вступить" />
      <GroupCard name="Группа 3" text="Описание группы 3" textButton="Вступить" />
    </div>
  );
};

export default MyGroups;
