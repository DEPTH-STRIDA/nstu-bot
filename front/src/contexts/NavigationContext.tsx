import React, { createContext, useContext, useState } from "react";

type NavigationContextType = {
  activePage: string;
  setActivePage: (page: string) => void;
};

const NavigationContext = createContext<NavigationContextType | undefined>(
  undefined
);

export const NavigationProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [activePage, setActivePage] = useState("поиск аудитории");

  return (
    <NavigationContext.Provider value={{ activePage, setActivePage }}>
      {children}
    </NavigationContext.Provider>
  );
};

export const useNavigation = () => {
  const context = useContext(NavigationContext);
  if (!context) {
    throw new Error("useNavigation must be used within NavigationProvider");
  }
  return context;
};
