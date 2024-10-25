import React, { useState } from 'react';
import './FileManager.css';

const FileManager = () => {
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [loading, setLoading] = useState(false);
  const [fileName, setFileName] = useState('');

  // Handle command execution
  const handleExecute = async () => {
    if (input.trim() !== '') {
      setLoading(true);
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
        setLoading(false);
      }
    } else {
      setOutput('Error: No se ingresó ningún comando');
    }
  };

  // Handle file upload
  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.name.endsWith('.smia')) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setInput(e.target.result);
        setFileName(file.name); // Set file name
      };
      reader.readAsText(file);
    } else {
      setOutput('Error: Solo se permiten archivos .smia');
      setFileName(''); // Clear file name if invalid
    }
  };

  return (
    <div className="gestor-archivos">
      <div className="panel-principal">
        <h1>Gestor de Archivos EXT2</h1>

        {/* Input section */}
        <div className="seccion-entrada">
          <label>Entrada:</label>
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ingresa el comando"
            disabled={loading}
          />
        </div>

        {/* Buttons section */}
        <div className="panel-botones">
          {/* Cargar Script Button */}
          <button className="boton-cargar" onClick={() => document.getElementById('file-upload').click()}>
            📂 Cargar Script
          </button>
          <input
            id="file-upload"
            type="file"
            accept=".smia"
            onChange={handleFileChange}
          />

          {/* Ejecutar Button */}
          <button className="boton-ejecutar" onClick={handleExecute} disabled={loading}>
            {loading ? '💻 Ejecutando...' : '💻 Ejecutar'}
          </button>
        </div>

        {/* File name display */}
        <div className="file-input-display">
          {fileName ? `📄 Archivo cargado: ${fileName}` : 'Ningún archivo cargado'}
        </div>

        {/* Output section */}
        <div className="seccion-salida">
          <label>Salida:</label>
          <textarea value={output} readOnly />
        </div>
      </div>
    </div>
  );
};

export default FileManager;
