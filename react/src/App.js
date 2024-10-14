import './App.css';
import DownloadForm from './components/download-form/DownloadForm';
import Version from './components/version/Version';


function App() {
  
  return (
    <div className="App">
      <header className="App-header">
          <DownloadForm />
          <Version />
      </header>
    </div>
  );
}

export default App;
