import Input from "./Input";
import { hint } from "../assets/img";
import { useHint } from "../contexts/HintContext";

interface GroupDetailsCardProps {
  name: string;
  description: string;
  onNameChange: (value: string) => void;
  onDescriptionChange: (value: string) => void;
}

const GroupDetailsCard = ({
  name,
  description,
  onNameChange,
  onDescriptionChange,
}: GroupDetailsCardProps) => {
  const { showHint } = useHint();

  const handleNameHintClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    showHint("Подсказка для названия группы");
  };

  return (
    <div className="bg-primary w-[320px] py-[35px] px-[20px] rounded-[15px] relative">
      <div
        className="absolute right-[7px] top-[4px] w-[25px] h-[25px] cursor-pointer
        before:content-[''] before:absolute before:-inset-3 before:cursor-pointer"
        onClick={handleNameHintClick}
      >
        <img src={hint} alt="Показать подсказку" className="w-full h-full" />
      </div>
      <div className="flex flex-col justify-start items-center w-full gap-[23px]">
        <Input
          placeholder="Название группы"
          value={name}
          onChange={onNameChange}
        />
        <Input
          placeholder="Описание группы"
          value={description}
          onChange={onDescriptionChange}
        />
      </div>
    </div>
  );
};

export default GroupDetailsCard;
