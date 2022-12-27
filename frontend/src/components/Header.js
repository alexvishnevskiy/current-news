import AppBar from "@mui/material/AppBar";
import logo from "../static/logo.jpeg";
import Toolbar from "@mui/material/Toolbar";
import Container from '@mui/material/Container';
import React from 'react';


export default function Header() {
    return (
        <AppBar position="static" color="default" style={{ background: '#cf473bf5' }}>
                <img src={logo} width="200" height="40" alt="current-news logo" className={"logo"} />
        </AppBar>
    )
}