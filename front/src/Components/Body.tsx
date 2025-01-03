import { useEffect, useState, useRef } from "react";
import { useNavigation } from "../contexts/NavigationContext";
import { SearchRoom, Main, SearchGroup, MyGroups, CreateGroup } from "../Pages";

/**
 * Основной компонент контента.
 * Управляет отображением страниц и свайп-навигацией между ними.
 * Синхронизируется с верхним меню для отображения активной страницы.
 */
const Body = () => {
  const { activePage, setActivePage } = useNavigation();
  const containerRef = useRef<HTMLDivElement>(null);
  const [startX, setStartX] = useState(0);
  const [startY, setStartY] = useState(0);
  const [currentX, setCurrentX] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [isScrollingVertically, setIsScrollingVertically] = useState(false);

  const menuItems = [
    "поиск аудитории",
    "главная",
    "поиск группы",
    "мои группы",
    "создать группу",
  ];

  const currentIndex = menuItems.indexOf(activePage);

  const handleTouchStart = (e: React.TouchEvent) => {
    setStartX(e.touches[0].clientX);
    setStartY(e.touches[0].clientY);
    setIsDragging(true);
    setIsScrollingVertically(false);
    if (containerRef.current) {
      containerRef.current.style.transition = "none";
    }
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging) return;

    const diffX = e.touches[0].clientX - startX;
    const diffY = e.touches[0].clientY - startY;

    // Определяем направление свайпа при начале движения
    if (!isScrollingVertically && Math.abs(diffY) > Math.abs(diffX)) {
      setIsScrollingVertically(true);
      return;
    }

    // Если определили, что это вертикальный скролл, игнорируем горизонтальное движение
    if (isScrollingVertically) return;

    // Предотвращаем вертикальный скролл при горизонтальном свайпе
    // e.preventDefault();

    const resistance = 0.3;
    let adjustedDiff = diffX;

    if (
      (currentIndex === 0 && diffX > 0) ||
      (currentIndex === menuItems.length - 1 && diffX < 0)
    ) {
      adjustedDiff = diffX * resistance;
    }

    setCurrentX(adjustedDiff);

    if (containerRef.current) {
      const baseOffset = -100 * currentIndex;
      containerRef.current.style.transform = `translateX(calc(${baseOffset}vw + ${adjustedDiff}px))`;
    }
  };

  const handleTouchEnd = () => {
    setIsDragging(false);
    if (containerRef.current) {
      containerRef.current.style.transition = "transform 0.3s ease-out";

      const threshold = window.innerWidth * 0.2; // 20% экрана для свайпа
      if (Math.abs(currentX) > threshold) {
        if (currentX > 0 && currentIndex > 0) {
          setActivePage(menuItems[currentIndex - 1]);
        } else if (currentX < 0 && currentIndex < menuItems.length - 1) {
          setActivePage(menuItems[currentIndex + 1]);
        } else {
          containerRef.current.style.transform = `translateX(-${
            currentIndex * 100
          }vw)`;
        }
      } else {
        containerRef.current.style.transform = `translateX(-${
          currentIndex * 100
        }vw)`;
      }
    }
    setCurrentX(0);
  };

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.style.transition = "transform 0.3s ease-out";
      containerRef.current.style.transform = `translateX(-${
        currentIndex * 100
      }vw)`;
    }
  }, [activePage]);

  useEffect(() => {
    const buttons = document.querySelectorAll("button");
    const activeButton = Array.from(buttons).find(
      (button) => button.textContent === activePage
    );

    if (activeButton) {
      activeButton.click();
    }
  }, [activePage]);

  return (
    <div className="fixed inset-0 mt-[140px] overflow-hidden">
      <div
        ref={containerRef}
        className="h-full flex"
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
        style={{
          willChange: "transform",
          transform: `translateX(-${currentIndex * 100}vw)`,
        }}
      >
        {menuItems.map((_, index) => (
          <div key={index} className="h-full min-w-full overflow-y-auto">
            {index === 0 && <SearchRoom />}
            {index === 1 && <Main />}
            {index === 2 && <SearchGroup />}
            {index === 3 && <MyGroups />}
            {index === 4 && <CreateGroup />}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Body;
