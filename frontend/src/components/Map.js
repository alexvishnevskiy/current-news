import React, { useState, useEffect} from 'react';
import WorldMap from "react-world-map";
import { Tooltip } from 'react-tooltip';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import axios from 'axios';
import 'react-tooltip/dist/react-tooltip.css';
import './map.css';


function renderCategories(categories) {
    return (
        <ol>
        {categories.map(topic => (
            <li key = {topic}>{topic}</li>
        ))}
        </ol>
    )
}

export default function Map() {
  const [dict, updataData] = useState({
    selected: null,
    continents: {},
  });

  useEffect(() => {
    axios.get('http://localhost:8888/data', {mode: 'no-cors'})
    .then(res => res.data)
    .then(res => {
        updataData({
          selected: false,
          continents: res
        });
      });
  }, []);

  const { selected, continents } = dict;

  return (
    <Grid
    container
    spacing={0}
    direction="column"
    alignItems="center"
    justify="center"
    style={{ minHeight: '100vh' }}
   >
    <Grid item xs={3}>
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
      }}
    >
    <Card sx={{ minWidth: 300}}>
        <CardContent>
            <Typography align="center" variant="h6" component="div">
                Ranking of article categories for each continent
            </Typography>
        <WorldMap selected={selected} onSelect={selected}/>
        <Tooltip anchorId='path4307' content={continents.Africa && renderCategories(continents.Africa)} events={['click']}/>
        <Tooltip anchorId='path5920' content={continents.Asia && renderCategories(continents.Asia)} events={['click']}/>
        <Tooltip anchorId='path5914' content={continents.Oceania && renderCategories(continents.Oceania)} events={['click']}/>
        <Tooltip anchorId='path5216' content={continents.Europe && renderCategories(continents.Europe)} events={['click']}/>
        <Tooltip anchorId='path5918' content={continents.SouthAmerica && renderCategories(continents.SouthAmerica)} events={['click']}/>
        <Tooltip anchorId='path5916' content={continents.NorthAmerica && renderCategories(continents.NorthAmerica)} events={['click']}/>
        </CardContent>
    </Card>
    </div>
    </Grid>
    </Grid>
  );
}