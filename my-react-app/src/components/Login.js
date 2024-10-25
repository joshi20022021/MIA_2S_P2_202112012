import React, { useState } from 'react';
import './Login.css'; // Importa los estilos del login

const Login = ({ onLoginSuccess }) => {
  const [usuario, setUsuario] = useState('');
  const [password, setPassword] = useState('');
  const [idParticion, setIdParticion] = useState(''); // Nuevo campo
  const [error, setError] = useState('');

  const handleLogin = (e) => {
    e.preventDefault();
    
    // Validación simple de usuario, contraseña e ID Partición
    if (usuario === 'root' && password === '123' && idParticion === '112A') {
      onLoginSuccess(); // Notificar al componente padre que el login fue exitoso
      setError('');
    } else {
      setError('Credenciales o ID incorrectos');
    }
  };

  return (
    <div className="login-container">
      <form onSubmit={handleLogin} className="login-form">
        <h2>Iniciar Sesión</h2>
        
        <div className="form-group">
          <label>Usuario</label>
          <input 
            type="text" 
            value={usuario} 
            onChange={(e) => setUsuario(e.target.value)} 
            placeholder="Usuario"
          />
        </div>

        <div className="form-group">
          <label>Contraseña</label>
          <input 
            type="password" 
            value={password} 
            onChange={(e) => setPassword(e.target.value)} 
            placeholder="Contraseña"
          />
        </div>

        {/* Nuevo campo ID Partición */}
        <div className="form-group">
          <label>ID Partición</label>
          <input 
            type="text" 
            value={idParticion} 
            onChange={(e) => setIdParticion(e.target.value)} 
            placeholder="ID Partición"
          />
        </div>
        
        {error && <p className="error">{error}</p>}

        <div className="form-group">
          <button type="submit">Ingresar</button>
        </div>
      </form>
    </div>
  );
};

export default Login;
