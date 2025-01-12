import React, { useState, useEffect, useRef } from "react";
import { menu, cross, delete_button,arrow_up,arrow_down } from "../assets/img";
import { createPortal } from "react-dom";

interface ScheduleRecordProps {
  index: number;
  isEven?: boolean;
  inputClassName?: string;
}

const ScheduleRecord = ({ index, isEven = false, inputClassName = "" }: ScheduleRecordProps) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node) &&
          buttonRef.current && !buttonRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const toggleMenu = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsMenuOpen(!isMenuOpen);
  };

  const renderMenu = () => {
    if (!buttonRef.current) return null;
    const rect = buttonRef.current.getBoundingClientRect();
    const top = rect.top + window.scrollY;
    const left = rect.left + window.scrollX;

    return createPortal(
      <div 
        ref={menuRef}
        className={`fixed p-[6px] gap-[3px]
          bg-inputBg/60 backdrop-blur-sm rounded-[20px]
          flex flex-col justify-center items-center
          transition-all duration-300 ease-in-out
          ${isMenuOpen 
            ? 'opacity-100 translate-x-0 scale-100' 
            : 'opacity-0 translate-x-2 scale-95 pointer-events-none'}`}
        style={{
          top: top + rect.height/2,
          left: left - 55,
          transform: `translateY(-50%) ${isMenuOpen ? 'translateX(0) scale(1)' : 'translateX(8px) scale(0.95)'}`,
          transformOrigin: 'left center'
        }}
      >
        <img src={arrow_up} alt="Вверх" className="w-[44px] h-[44px] cursor-pointer" />
        <img src={delete_button} alt="Удалить" className="w-[44px] h-[44px] cursor-pointer" />
        <img src={arrow_down} alt="Вниз" className="w-[44px] h-[44px] cursor-pointer" />
      </div>,
      document.body
    );
  };

  return (
    <div
      className={`w-full grid grid-cols-[30px_55px_1fr_40px_50px] items-center relative
      ${isEven ? "bg-weekDayBg" : "bg-white"}
      divide-x divide-primary border-b border-[0.5px] border-primary text-primary font-roboto font-light`}
    >
      <div className="flex justify-center items-center h-[55px] text-[13px]">
        {index}
      </div>
      <div className="flex justify-center items-center h-[55px]">
        <textarea
          placeholder="Время"
          className={`w-full h-full bg-transparent text-[13px] resize-none 
            text-center py-[18px] leading-[18px]
            text-primary placeholder:text-primary focus:outline-none ${inputClassName}`}
        />
      </div>
      <div className="flex justify-center items-center h-[55px]">
        <textarea
          placeholder="Название предмета"
          className={`w-full h-full bg-transparent text-[13px] resize-none
          text-center py-[18px] leading-[18px]
            text-primary placeholder:text-primary focus:outline-none ${inputClassName}`}
        />
      </div>
      <div className="flex justify-center items-center h-[55px]">
        <textarea
          placeholder="№"
          className={`w-full h-full bg-transparent text-[13px] resize-none
           text-center py-[18px] leading-[18px]
            text-primary placeholder:text-primary focus:outline-none ${inputClassName}`}
        />
      </div>
      <div ref={buttonRef} className="flex justify-center items-center h-[55px] relative">
        <img 
          src={isMenuOpen ? cross : menu} 
          alt={isMenuOpen ? "Закрыть меню" : "Открыть меню"} 
          className={`w-[43px] h-[43px] cursor-pointer transition-all duration-300
          hover:scale-110 transform ${isMenuOpen ? 'scale-110' : 'scale-100'}`}
          onClick={toggleMenu}
        />
        {renderMenu()}
      </div>
    </div>
  );
};

export default ScheduleRecord; 