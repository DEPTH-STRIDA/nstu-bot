import SearchBar from "../Components/SearchBar";

/**
 * Страница поиска группы.
 * Предоставляет интерфейс для поиска учебных групп.
 */
const SearchGroup = () => {
  return (
    <div className=" w-full bg-white flex flex-col justify-start items-center pt-[15px]">
      <SearchBar placeholder="Введите название группы" onChange={() => {}} />
    </div>
  );
};

export default SearchGroup;
