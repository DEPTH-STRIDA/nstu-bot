import { hamburger } from "../assets/img";

interface HeaderProps {
  onMenuClick: () => void;
}

const Header = ({ onMenuClick }: HeaderProps) => {
  return (
    <div className="flex flex-col justify-center items-center pt-[18px] pb-[15px] relative">
      <img
        src={hamburger}
        alt="open menu"
        className="absolute left-[18px] top-[23px] cursor-pointer"
        onClick={onMenuClick}
      />
      <h2 className="text-[13px] font-light">
        <span className="font-bold">GOLANG</span> DEVELOPER
      </h2>
      <div className="w-[100px] h-[1px] bg-secondary mt-[8px] mb-[4px]"></div>
      <h2 className="text-2xl font-light">
        <span className="font-bold">SHARAGA</span> BOT
      </h2>
    </div>
  );
};

export default Header;
