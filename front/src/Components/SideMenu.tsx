import Menu from "./Menu";

interface SideMenuProps {
  isOpen: boolean;
  onClose: () => void;
  menuItems: string[];
  activePage: string;
  onItemClick: (item: string) => void;
}

/**
 * Компонент бокового меню.
 * Отображает выдвижное меню с затемнением фона.
 * Содержит те же пункты, что и верхнее меню.
 */
const SideMenu = ({
  isOpen,
  onClose,
  menuItems,
  activePage,
  onItemClick,
}: SideMenuProps) => {
  return (
    <div
      className={`fixed inset-0 bg-black/50 backdrop-blur-sm transition-opacity z-[60] ${
        isOpen ? "opacity-100" : "opacity-0 pointer-events-none"
      }`}
      onClick={onClose}
    >
      <div
        className={`fixed left-0 top-0 h-full w-[250px]
           bg-primary
           transform transition-transform duration-300 ${
             isOpen ? "translate-x-0" : "-translate-x-full"
           }`}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="pt-[45px] pl-[15px]">
          <Menu
            items={menuItems}
            activePage={activePage}
            onItemClick={onItemClick}
            variant="side"
          />
        </div>
      </div>
    </div>
  );
};

export default SideMenu;
