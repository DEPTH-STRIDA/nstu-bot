import { hint } from "../assets/img";
import { useNavigation } from "../contexts/NavigationContext";
import { useHint } from "../contexts/HintContext";

interface InfoCardProps {
  text: string;
  hintText?: string;
  targetPage: string;
}

/**
 * Информационная карточка для главной страницы
 * @param {string} text - Основной текст карточки
 * @param {string} hintText - Дополнительный текст-подсказка (опционально)
 * @param {string} targetPage - Целевая страница для навигации
 */
const InfoCard = ({ text, hintText, targetPage }: InfoCardProps) => {
  const { setActivePage } = useNavigation();
  const { showHint } = useHint();

  const handleCardClick = () => {
    setActivePage(targetPage);
  };

  const handleHintClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    showHint(hintText || "");
  };

  return (
    <div
      className="
    flex flex-col justify-center items-center
    bg-primary rounded-[15px] py-[30px] w-[300px]
    relative cursor-pointer"
      onClick={handleCardClick}
    >
      <div className="text-secondary text-[20px] font-light font-roboto">
        {text}
      </div>
      {hintText && (
        <div
          className="absolute right-[7px] top-[4px] w-[25px] h-[25px] cursor-pointer
          before:content-[''] before:absolute before:-inset-3 before:cursor-pointer"
          onClick={handleHintClick}
        >
          <img src={hint} alt="Показать подсказку" className="w-full h-full" />
        </div>
      )}
    </div>
  );
};

export default InfoCard;
