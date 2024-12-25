import { useHint } from "../contexts/HintContext";
import { close } from "../assets/img";

const MainHint = () => {
  const { isVisible, hideHint, hintText } = useHint();

  if (!isVisible) return null;

  return (
    <div className="">
      <div
        className="fixed inset-0 bg-white/5 backdrop-blur-sm z-[999] touch-none 
        animate-fadeIn"
        onClick={hideHint}
      />
      <div
        className="fixed top-[30%] left-1/2 -translate-x-1/2 -translate-y-1/2 bg-primary
        pt-[40px] pb-[30px] px-[20px] w-[350px]
        rounded-[15px] shadow-xl z-[1000] text-secondary text-[20px] font-light touch-none
        animate-popIn"
        onClick={(e) => e.stopPropagation()}
      >
        <div 
          className="absolute top-2 right-2 w-[30px] h-[30px] cursor-pointer
          before:content-[''] before:absolute before:-inset-3 before:cursor-pointer"
          onClick={hideHint}
        >
          <img
            src={close}
            alt="Закрыть"
            className="w-full h-full hover:opacity-80 transition-opacity"
          />
        </div>
        {hintText}
      </div>
    </div>
  );
};

export default MainHint;
