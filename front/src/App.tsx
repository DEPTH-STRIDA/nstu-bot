import Nav from "./Components/Nav";
import Body from "./Components/Body";
import MainHint from "./MainPage/MainHint";
import { NavigationProvider } from "./contexts/NavigationContext";
import { HintProvider } from "./contexts/HintContext";

const App = () => {
  return (
    <NavigationProvider>
      <HintProvider>
        <Nav />
        <Body />
        <MainHint />
      </HintProvider>
    </NavigationProvider>
  );
};

export default App;
