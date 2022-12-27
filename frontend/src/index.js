import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import Map from './components/Map'
import VirtualizedList from './components/Headlines'
import Header from './components/Header';
import Footer from './components/Footer'
import './index.css'

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <Header/>
        <App/>
        <Footer/>
    </React.StrictMode>
);
