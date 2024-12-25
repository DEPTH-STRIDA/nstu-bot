import { hint } from "../assets/img";
import { useHint } from "../contexts/HintContext";

interface InfoCardProps {
  text: string;
  hintText?: string;
}

/**
 * Информационная карточка для главной страницы
 * @param {string} text - Основной текст карточки
 * @param {string} hintText - Дополнительный текст-подсказка (опционально)
 */
const InfoCard = ({ text, hintText }: InfoCardProps) => {
  const { showHint } = useHint();

  return (
    <div
      className="
    flex flex-col justify-center items-center
    bg-primary rounded-[15px] py-[30px] w-[300px]
    relative
   "
    >
      <div className="text-secondary text-[20px] font-light font-roboto">
        {text}
      </div>
      {hintText && (
        <div
          className="absolute right-[7px] top-[4px] w-[25px] h-[25px] cursor-pointer
          before:content-[''] before:absolute before:-inset-3 before:cursor-pointer"
          onClick={(e) => {
            e.stopPropagation();
            showHint(hintText);
          }}
        >
          <img src={hint} alt="Показать подсказку" className="w-full h-full" />
        </div>
      )}
    </div>
  );
};

export default InfoCard;
