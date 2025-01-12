import { useState, useEffect } from "react";
import { hint } from "../assets/img";
import { useHint } from "../contexts/HintContext";

interface AlternationCardProps {
  onRetrieveCheckboxState: (checked: boolean) => void;
  onRetrieveDate: (date: string) => void;
}

const AlternationCard = ({
  onRetrieveCheckboxState,
  onRetrieveDate,
}: AlternationCardProps) => {
  const { showHint } = useHint();
  const [isChecked, setIsChecked] = useState(false);
  const [date, setDate] = useState("");

  useEffect(() => {
    onRetrieveCheckboxState(isChecked);
  }, [isChecked, onRetrieveCheckboxState]);

  useEffect(() => {
    onRetrieveDate(date);
  }, [date, onRetrieveDate]);

  const handleAlternationClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    showHint(
      "При включенном чередовании недель необходимо указать первую нечетную неделю. Система будет вести отсчет, начиная с этой недели."
    );
  };

  const handleCheckboxChange = () => {
    setIsChecked(!isChecked);
  };

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setDate(e.target.value);
  };

  return (
    <div className="bg-primary w-[320px] py-[35px] px-[20px] rounded-[15px] relative mt-[24px]">
      <div
        className="absolute right-[7px] top-[4px] w-[25px] h-[25px] cursor-pointer
        before:content-[''] before:absolute before:-inset-3 before:cursor-pointer"
        onClick={handleAlternationClick}
      >
        <img src={hint} alt="Показать подсказку" className="w-full h-full" />
      </div>

      <div
        className="flex flex-row justify-center items-center w-full gap-[18px]  
      text-secondary text-[20px] font-light font-roboto"
      >
        <input
          type="checkbox"
          checked={isChecked}
          onChange={handleCheckboxChange}
          className="form-checkbox h-[30px] w-[30px] text-primary transition duration-150 ease-in-out"
        />
        Чредование недель
      </div>

      <div
        className={`overflow-hidden transition-max-height duration-500 ease-in-out ${
          isChecked ? "max-h-[200px]" : "max-h-0"
        }`}
      >
        <div className="flex flex-col justify-center items-center w-full gap-[14px] px-[20px] opacity-100">
          <p className="font-roboto text-[20px] font-light text-secondary mt-[26px]">
            Укажите понедельник любой нечетной недели
          </p>
          <input
            type="date"
            name="date"
            id="date"
            value={date}
            onChange={handleDateChange}
            className="min-w-[220px] h-[40px] rounded-[15px] 
            flex justify-center items-center font-roboto text-[24px] font-light text-primary"
          />
        </div>
      </div>
    </div>
  );
};

export default AlternationCard;
