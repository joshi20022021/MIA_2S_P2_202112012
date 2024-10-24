import React, { useState } from 'react';
import './FileManager.css';

const FileManager = () => {
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [loading, setLoading] = useState(false);  // Estado de carga

  // Maneja la ejecución del comando
  const handleExecute = async () => {
    if (input.trim() !== '') {
      setLoading(true);  // Inicia el estado de carga
      try {
        const response = await fetch('http://localhost:8080/api/command', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ command: input }),
        });

        const data = await response.json();
        setOutput(data.output);
      } catch (error) {
        setOutput('Error al comunicarse con el servidor');
      } finally {
        setLoading(false);  // Finaliza el estado de carga
      }
    } else {
      setOutput('Error: No se ingresó ningún comando');
    }
  };

  // Maneja la carga del archivo .smia
  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.name.endsWith('.smia')) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setInput(e.target.result);
      };
      reader.readAsText(file);
    } else {
      setOutput('Error: Solo se permiten archivos .smia');
    }
  };

  return (
    <div className="gestor-archivos">
      {/* Panel lateral para los botones */}
      <div className="panel-lateral">
        <label htmlFor="file-upload" className="boton-cargar">
          📂 Cargar Script
        </label>
        <input
          id="file-upload"
          type="file"
          accept=".smia"
          style={{ display: 'none' }}
          onChange={handleFileChange}
        />
        <button className="boton-ejecutar" onClick={handleExecute} disabled={loading}>
          {loading ? '💻 Ejecutando...' : '💻 Ejecutar'}
        </button>
      </div>

      {/* Panel principal con entrada y salida */}
      <div className="panel-principal">
        <h1>Gestor de Archivos EXT2</h1>
        <div className="seccion-entrada">
          <label>Entrada:</label>
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ingresa el comando"
            disabled={loading}  // Deshabilita mientras está cargando
          />
        </div>
        <div className="seccion-salida">
          <label>Salida:</label>
          <textarea value={output} readOnly />
        </div>
      </div>
    </div>
  );
};

export default FileManager;
