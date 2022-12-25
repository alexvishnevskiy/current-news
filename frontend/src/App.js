import React from 'react';
import Headlines from './components/Headlines';
import Map from './components/Map';
import Grid from '@mui/material/Grid';

function App() {
  return (
    <Grid container spacing={4} direction="row">
      <Grid item md={7}>
        <Map/>
      </Grid>
      <Grid item md={5}>
        <Headlines/>
      </Grid>
    </Grid>
  )
}

export default App;
