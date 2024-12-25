import { createContext, useContext, useState, ReactNode } from "react";

interface HintContextType {
  showHint: (text: string) => void;
  hideHint: () => void;
  hintText: string | null;
  isVisible: boolean;
}

const HintContext = createContext<HintContextType | null>(null);

export const useHint = () => {
  const context = useContext(HintContext);
  if (!context) {
    throw new Error("useHint must be used within a HintProvider");
  }
  return context;
};

interface HintProviderProps {
  children: ReactNode;
}

export const HintProvider = ({ children }: HintProviderProps) => {
  const [hintText, setHintText] = useState<string | null>(null);
  const [isVisible, setIsVisible] = useState(false);

  const showHint = (text: string) => {
    setHintText(text);
    setIsVisible(true);
  };

  const hideHint = () => {
    setIsVisible(false);
    setHintText(null);
  };

  return (
    <HintContext.Provider value={{ showHint, hideHint, hintText, isVisible }}>
      {children}
    </HintContext.Provider>
  );
}; 