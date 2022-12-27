import React from 'react';
import Headlines from './components/Headlines';
import ArchiveChart from './components/Chart'
import Map from './components/Map';
import Grid from '@mui/material/Grid';

function App() {
  return (
    <Grid container spacing={4} direction="row" sx = {{paddingRight: 2, paddingLeft: 2, paddingTop: 3}}>
      <Grid item xs={12} md={12}>
        <ArchiveChart/>
      </Grid>
      <Grid item xs={12} md={7}>
        <Map/>
      </Grid>
      <Grid item xs={12} md={5}>
        <Headlines/>
      </Grid>
    </Grid>
  )
}

export default App;
