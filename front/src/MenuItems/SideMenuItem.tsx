interface SideMenuItemProps {
  item: string;
  isActive: boolean;
  onClick: (item: string) => void;
}

/**
 * Компонент пункта бокового меню
 */
const SideMenuItem = ({ item, isActive, onClick }: SideMenuItemProps) => {
  return (
    <button
      className={`
        w-[200px] 
        text-left text-secondary
        font-roboto font-light
        text-[18px]
        px-[13px] py-[4px] rounded-[8px]
        flex items-center
        transition-all duration-300
        ${isActive ? "bg-secondary/20" : "hover:bg-secondary/10"}
      `}
      onClick={() => onClick(item)}
    >
      {item}
    </button>
  );
};

export default SideMenuItem;
