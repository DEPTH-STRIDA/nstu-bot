import TopMenuItem from "../MenuItems/TopMenuItem";
import SideMenuItem from "../MenuItems/SideMenuItem";

interface MenuProps {
  items: string[];
  activePage: string;
  onItemClick: (item: string) => void;
  variant: "top" | "side";
  onButtonClick?: (
    e: React.MouseEvent<HTMLButtonElement>,
    item: string
  ) => void;
}

/**
 * Компонент меню, который может быть использован как для верхнего, так и для бокового меню.
 * Отображает список пунктов меню с соответствующими стилями.
 * @param {string[]} items - Массив пунктов меню
 * @param {string} activePage - Текущая активная страница
 * @param {function} onItemClick - Обработчик клика по пункту меню
 * @param {'top' | 'side'} variant - Вариант отображения (верхнее или боковое меню)
 * @param {function} onButtonClick - Дополнительный обработчик для верхнего меню
 */
const Menu = ({
  items,
  activePage,
  onItemClick,
  variant,
  onButtonClick,
}: MenuProps) => {
  const handleClick = (
    e: React.MouseEvent<HTMLButtonElement>,
    item: string
  ) => {
    if (onButtonClick) {
      onButtonClick(e, item);
    } else {
      onItemClick(item);
    }
  };

  const MenuItem = variant === "top" ? TopMenuItem : SideMenuItem;

  return (
    <div className={variant === "top" ? "flex gap-6" : "flex flex-col gap-4"}>
      {items.map((item) => (
        <div key={item} onClick={(e) => handleClick(e as any, item)}>
          <MenuItem
            item={item}
            isActive={activePage === item}
            onClick={() => {}}
          />
        </div>
      ))}
    </div>
  );
};

export default Menu;
