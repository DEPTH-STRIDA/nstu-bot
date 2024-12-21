import Nav from "./Components/Nav";
import Body from "./Components/Body";
import { NavigationProvider } from "./contexts/NavigationContext";

const App = () => {
  return (
    <NavigationProvider>
      <Nav />
      <Body />
    </NavigationProvider>
  );
};

export default App;
