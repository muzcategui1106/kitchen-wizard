import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import { TopNavigationBar } from './Navigation/TopBar';

export default function App() {
  return (
    <div className="App">
      <div className="Navigation-Bar">
        <TopNavigationBar/>
      </div>
    </div>
  );
}
