import { seach } from "../assets/img";

interface SearchBarProps {
  placeholder: string;
  onChange: (value: string) => void;
}

/**
 * Компонент поисковой строки
 * @param {string} placeholder - Текст-подсказка в поле ввода
 * @param {function} onChange - Функция получения введенного значения
 */
const SearchBar = ({ placeholder, onChange }: SearchBarProps) => {
  return (
    <div
      className="flex flex-row justify-start items-center bg-inputBg
    h-[40px] w-[350px] rounded-[15px] px-[13px] mb-[26px]"
    >
      <img
        src={seach}
        alt="Иконка лупы"
        className="
      w-[30px] pr-[13px]"
      />
      <input
        type="text"
        onChange={(e) => onChange(e.target.value)}
        className="bg-inputBg w-full
        text-[16px] font-normal  text-primary
        placeholder:text-primary/50
        focus:outline-none"
        placeholder={placeholder}
      />
    </div>
  );
};

export default SearchBar;
