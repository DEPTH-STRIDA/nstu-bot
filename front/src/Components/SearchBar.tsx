import { useState } from "react";
import { seach } from "../assets/img";

interface SearchBarProps {
  placeholder: string;
  onChange: (value: string) => void;
  smallMargin?: boolean;
}

/**
 * Компонент поисковой строки
 * @param {string} placeholder - Текст-подсказка в поле ввода
 * @param {function} onChange - Функция получения введенного значения
 * @param {boolean} smallMargin - Уменьшенный отступ сверху
 */
const SearchBar = ({ placeholder, onChange, smallMargin }: SearchBarProps) => {
  const [isFocused, setIsFocused] = useState(false);
  const [value, setValue] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setValue(newValue);
    onChange(newValue);
  };

  return (
    <div
      className={`flex flex-row justify-start items-center bg-inputBg
    h-[40px] w-[350px] rounded-[15px] px-[13px] ${
      smallMargin ? "mb-[19px]" : "mb-[26px]"
    } relative`}
    >
      <img
        src={seach}
        alt="Иконка лупы"
        className="
      w-[30px] pr-[13px]"
      />
      <input
        type="text"
        value={value}
        onChange={handleChange}
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        className="bg-inputBg w-full
        text-[16px] font-normal text-primary
        focus:outline-none"
      />
      {!isFocused && !value && (
        <div
          className="absolute inset-0 flex justify-center items-center text-primary/50 pointer-events-none
        text-[24px] font-light font-roboto"
        >
          {placeholder}
        </div>
      )}
    </div>
  );
};

export default SearchBar;
