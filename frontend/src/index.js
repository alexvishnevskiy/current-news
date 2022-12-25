import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import Map from './components/Map'
import VirtualizedList from './components/Headlines'
import './index.css'

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
);
