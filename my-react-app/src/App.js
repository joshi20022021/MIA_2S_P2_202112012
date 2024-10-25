import React, { useState } from 'react';
import FileManager from './components/FileManager'; // Importa el FileManager
import Login from './components/Login'; // Importa el Login

const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false); // Estado de login

  // Función que se ejecuta cuando el login es exitoso
  const handleLoginSuccess = () => {
    setIsLoggedIn(true);
  };

  return (
    <div>
      {isLoggedIn ? <FileManager /> : <Login onLoginSuccess={handleLoginSuccess} />}
    </div>
  );
};

export default App;
