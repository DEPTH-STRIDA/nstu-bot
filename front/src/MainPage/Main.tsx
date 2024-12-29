import InfoCard from "./InfoCard";

/**
 * Главная страница приложения.
 * Отображает основную информацию и функционал.
 */
const Main = () => {
  return (
    <div
      className="w-full bg-white pt-[40px]
    flex flex-col justify-start items-center text-center"
    >
      <div className="flex flex-col justify-center items-center gap-[18px]">
        <InfoCard
          text="Присоединиться к группе"
          hintText="Нажмите, чтобы присоединиться к существующей группе"
          targetPage="поиск группы"
        />
        <InfoCard
          text="Создать группу"
          hintText="Выберите нужный раздел в меню"
          targetPage="создать группу"
        />
        <InfoCard
          text="Мои группы"
          hintText="Выберите нужный раздел в меню"
          targetPage="мои группы"
        />
        <InfoCard
          text="Поиск аудитории"
          hintText="Выберите нужный раздел в меню"
          targetPage="поиск аудитории"
        />
      </div>
    </div>
  );
};

export default Main;
