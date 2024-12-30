import { ChangeEvent, useState } from "react";

interface InputProps {
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}

/**
 * Компонент ввода
 * @param {string} placeholder - Текст-подсказка в поле ввода
 * @param {string} value - Значение инпута
 * @param {function} onChange - Функция получения введенного значения
 */
const Input = ({ placeholder, value, onChange }: InputProps) => {
  const [isFocused, setIsFocused] = useState(false);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value);
  };

  return (
    <div
      className={`flex flex-row justify-start items-center bg-inputBg
    h-[40px] w-full rounded-[15px] px-[13px] relative`}
    >
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
      {!value && !isFocused && (
        <div
          className="absolute inset-0 flex justify-center items-center text-primary pointer-events-none
        text-[24px] font-light font-roboto"
        >
          {placeholder}
        </div>
      )}
    </div>
  );
};

export default Input;
