import { useRef, useState, useEffect } from "react";
import { useNavigation } from "../contexts/NavigationContext";
import { hamburger } from "../assets/img";

const Nav = () => {
  const { activePage, setActivePage } = useNavigation();
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [activeButtonWidth, setActiveButtonWidth] = useState(0);
  const [activeButtonLeft, setActiveButtonLeft] = useState(0);

  const menuItems = [
    "поиск аудитории",
    "главная",
    "поиск группы",
    "ваши группы",
    "создать группу",
  ];

  const scrollToActiveButton = (button: HTMLButtonElement) => {
    if (scrollContainerRef.current) {
      const container = scrollContainerRef.current;
      const scrollLeft =
        button.offsetLeft - container.offsetWidth / 2 + button.offsetWidth / 2;
      container.scrollTo({ left: scrollLeft, behavior: "smooth" });

      setActiveButtonWidth(button.offsetWidth);
      setActiveButtonLeft(button.offsetLeft);
    }
  };

  useEffect(() => {
    const buttons = document.querySelectorAll("button");
    const activeButton = Array.from(buttons).find(
      (button) => button.textContent === activePage
    );

    if (activeButton) {
      setActiveButtonWidth(activeButton.offsetWidth);
      setActiveButtonLeft(activeButton.offsetLeft);
    }
  }, [activePage]);

  return (
    <div className="bg-primary text-secondary w-full flex-col justify-center items-center fixed top-0 z-50">
      <div className="flex flex-col justify-center items-center pt-[18px] pb-[15px] relative">
        <img
          src={hamburger}
          alt="open menu"
          className="absolute left-[18px] top-[23px]"
        />
        <h2 className="text-[13px] font-light">
          <span className="font-bold">GOLANG</span> DEVELOPER
        </h2>
        <div className="w-[100px] h-[1px] bg-secondary mt-[8px] mb-[4px]"></div>
        <h2 className="text-2xl font-light">
          <span className="font-bold">SHARAGA</span> BOT
        </h2>
      </div>
      <div
        className="w-full overflow-x-auto scrollbar-hide"
        ref={scrollContainerRef}
      >
        <div className="scrollbar-hide flex flex-row justify-between items-center px-[18px] font-light text-[14px] font-roboto pb-[15px] min-w-max gap-6 relative">
          {menuItems.map((item) => (
            <button
              key={item}
              className="relative py-1"
              onClick={(e) => {
                setActivePage(item);
                scrollToActiveButton(e.currentTarget);
              }}
            >
              {item}
            </button>
          ))}
          <div
            className="absolute bottom-[11px] h-[2px] bg-secondary shadow-activeMenu transition-all duration-300 ease-in-out"
            style={{
              left: `${activeButtonLeft}px`,
              width: `${activeButtonWidth}px`,
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default Nav;
