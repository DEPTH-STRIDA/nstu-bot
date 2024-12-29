import { consists } from "../assets/img";

interface GroupCardProps {
  name: string;
  text: string;
  textButton: string;
  consist?: boolean;
}

const GroupCard: React.FC<GroupCardProps> = ({ name, text, textButton, consist }) => {
  const truncatedName = name.length > 50 ? name.slice(0, 50) + '...' : name;
  const truncatedText = text.length > 110 ? text.slice(0, 110) + '...' : text;

  return (
    <div className="w-full max-w-[264px] bg-primary rounded-[15px] p-3 text-center relative">
      {consist && (
        <div className="absolute right-2 top-2 w-[22px] h-[22px]">
          <img
            src={consists}
            alt="Состоит в группе"
            className="w-full h-full"
          />
        </div>
      )}
      <h2 className="text-white text-[24px] font-roboto ">
        {truncatedName}
      </h2>
      <div className=" w-[137px] h-[1px] bg-secondary  mb-[3px] mx-auto"></div>
      <p className="text-white text-[17px] mb-1 font-roboto font-light">
        {truncatedText}
      </p>
      <button className="w-[60%] leading-[12.5px] bg-white text-black text-[16px] font-roboto font-light py-2 rounded-md hover:bg-gray-200 transition-colors">
        {textButton}
      </button>
    </div>
  );
};

export default GroupCard;
