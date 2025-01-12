import { useState, useEffect, useRef } from "react";
import { useNavigation } from "../contexts/NavigationContext";
import { hamburger } from "../assets/img";
import Menu from "./Menu";
import SideMenu from "./SideMenu";

/**
 * Главный компонент навигации.
 * Содержит верхнее меню с кнопками и боковое выдвижное меню.
 * Управляет состоянием активной страницы и анимацией подчеркивания активного пункта.
 */
const Nav = () => {
  const { activePage, setActivePage } = useNavigation();
  const [activeButtonWidth, setActiveButtonWidth] = useState(0);
  const [activeButtonLeft, setActiveButtonLeft] = useState(0);
  const [isSideMenuOpen, setIsSideMenuOpen] = useState(false);
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  const menuItems = [
    "поиск аудитории",
    "главная",
    "поиск группы",
    "мои группы",
    "создать группу",
  ];

  const scrollToActiveButton = (button: HTMLButtonElement) => {
    if (scrollContainerRef.current) {
      const container = scrollContainerRef.current;
      const scrollLeft =
        button.offsetLeft - container.offsetWidth / 2 + button.offsetWidth / 2;
      container.scrollTo({ left: scrollLeft, behavior: "smooth" });
    }
  };

  const handleButtonClick = (e: React.MouseEvent<HTMLButtonElement>, item: string) => {
    const button = e.currentTarget;
    setActiveButtonWidth(button.offsetWidth);
    setActiveButtonLeft(button.offsetLeft);
    setActivePage(item);
    scrollToActiveButton(button);
  };

  const handleSideMenuClick = (item: string) => {
    setActivePage(item);
    setIsSideMenuOpen(false);
  };

  useEffect(() => {
    const buttons = document.querySelectorAll("button");
    const activeButton = Array.from(buttons).find(
      (button) => button.textContent === activePage
    );

    if (activeButton) {
      setActiveButtonWidth(activeButton.offsetWidth);
      setActiveButtonLeft(activeButton.offsetLeft);
      scrollToActiveButton(activeButton);
    }
  }, [activePage]);

  return (
    <>
      <div className="bg-primary text-secondary w-full flex-col justify-center items-center fixed top-0 z-50">
        <div className="flex flex-col justify-center items-center pt-[18px] pb-[15px] relative">
          <img
            src={hamburger}
            alt="open menu"
            className="absolute left-[18px] top-[23px] cursor-pointer"
            onClick={() => setIsSideMenuOpen(true)}
          />
          <h2 className="text-[13px] font-light">
            <span className="font-bold">GOLANG</span> DEVELOPER
          </h2>
          <div className="w-[100px] h-[1px] bg-secondary mt-[8px] mb-[4px]"></div>
          <h2 className="text-2xl font-light">
            <span className="font-bold">SHARAGA</span> BOT
          </h2>
        </div>

        <div className="w-full overflow-x-auto scrollbar-hide" ref={scrollContainerRef}>
          <div className="scrollbar-hide px-[18px] font-light text-[14px] font-roboto pb-[15px] min-w-max relative">
            <Menu
              items={menuItems}
              activePage={activePage}
              onItemClick={() => {}}
              variant="top"
              onButtonClick={handleButtonClick}
            />
            <div
              className="absolute bottom-[16px] h-[2px] bg-secondary shadow-activeMenu transition-all duration-300 ease-in-out"
              style={{
                left: `${activeButtonLeft}px`,
                width: `${activeButtonWidth}px`,
              }}
            />
          </div>
        </div>
      </div>

      <SideMenu
        isOpen={isSideMenuOpen}
        onClose={() => setIsSideMenuOpen(false)}
        menuItems={menuItems}
        activePage={activePage}
        onItemClick={handleSideMenuClick}
      />
    </>
  );
};

export default Nav;
