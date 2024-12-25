interface TopMenuItemProps {
  item: string;
  isActive: boolean;
  onClick: (item: string) => void;
}

/**
 * Компонент пункта верхнего меню
 */
const TopMenuItem = ({ item, onClick }: TopMenuItemProps) => {
  return (
    <button
      className="relative py-1 text-secondary text-[14px] font-light"
      onClick={() => onClick(item)}
    >
      {item}
    </button>
  );
};

export default TopMenuItem; 