import React, { useState, useEffect} from 'react';
import WorldMap from "./WorldMap";
import { Tooltip } from 'react-tooltip';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import axios from 'axios';
import 'react-tooltip/dist/react-tooltip.css';
import './map.css';
import LinearProgress from '@mui/material/LinearProgress';
import NorthAmericaImage from "../static/images/NorthAmerica.jpg";
import AfricaImage from "../static/images/Africa.jpeg"
import AsiaImage from "../static/images/Asia.jpeg"
import EuropeImage from "../static/images/Europe.png"
import OceaniaImage from "../static/images/Oceania.jpeg"
import SouthAmericaImage from "../static/images/SouthAmerica.svg"


function renderCategories(categories) {
    return (
        <ol>
        {categories.map(topic => (
            <li key = {topic}>{topic}</li>
        ))}
        </ol>
    )
}

function GetImage(continent) {
  switch(continent){
    case 'Africa':
      return AfricaImage
    case 'Asia':
      return AsiaImage
    case 'SouthAmerica':
      return SouthAmericaImage
    case 'NorthAmerica':
      return NorthAmericaImage
    case 'Oceania':
      return OceaniaImage
    case 'Europe':
      return EuropeImage
    default:
      break;
  }
}

export default function Map() {
  const [dict, updataData] = useState({
    selected: null,
    continents: {},
  });

  useEffect(() => {
    axios.get('http://localhost:8888/categories', {mode: 'no-cors'})
    // axios.get('http://4.227.248.158/categories', {mode: 'no-cors'})
    .then(res => res.data)
    .then(res => {
        updataData({
          selected: false,
          continents: res
        });
      });
  }, []);

  const { selected, continents } = dict;
  const total = Object.entries(continents)
                .map(([key, value]) => value.Total)
                .reduce((accumulator, currentValue) => accumulator + currentValue, 0);
  // Create items array
  var items = Object.keys(continents).map(function(key) {
    return [key, continents[key]];
  });

  // Sort the array based on the second element
  items.sort(function(first, second) {
    return second[1].Total - first[1].Total;
  });

  return (
    <Card>
      <Grid container spacing={1} direction="row" alignItems="center" justify="center">
        <Grid item md={7}>
          <CardContent>
              <Typography align="left" variant="h7" component="div">
                <b>Ranking of article categories for each continent</b>
              </Typography>
          <WorldMap selected={selected} onSelect={selected}/>
          <Tooltip anchorId='path4307' content={continents.Africa && renderCategories(continents.Africa.Categories)} events={['click']}/>
          <Tooltip anchorId='path5920' content={continents.Asia && renderCategories(continents.Asia.Categories)} events={['click']}/>
          <Tooltip anchorId='path5914' content={continents.Oceania && renderCategories(continents.Oceania.Categories)} events={['click']}/>
          <Tooltip anchorId='path5216' content={continents.Europe && renderCategories(continents.Europe.Categories)} events={['click']}/>
          <Tooltip anchorId='path5918' content={continents.SouthAmerica && renderCategories(continents.SouthAmerica.Categories)} events={['click']}/>
          <Tooltip anchorId='path5916' content={continents.NorthAmerica && renderCategories(continents.NorthAmerica.Categories)} events={['click']}/>
          </CardContent>
        </Grid>
        <Grid item md={5}>
          <CardContent>
            <Box pl={3}>
              <Typography padding-left={10} align="left" variant="h7" component="div">
                <b>Leaderboard</b>
              </Typography>
            </Box>
            <List
              sx={{
                width: '100%',
                maxHeight: 450,
                maxWidth: 300,
                bgcolor: 'background.paper',
              }}
              component="div" disablePadding
            >
              {
                <React.Fragment>
                {items.map(([key, value]) => // cons is the array fetched from API that represents the ranking
                    <ListItem key={key} sx={{ height: 51 }}>
                      <ListItemAvatar>
                        <Avatar src={GetImage(key)} variant="square"/>
                      </ListItemAvatar>
                      <Box sx={{ width: '100%' }}>
                        <Stack direction="row" spacing={2}>
                          <ListItemText primary={`${key}`}/>
                          <Box mt={2}>
                            <ListItemText primary={`${total == 0 ? 0 : Math.round(value.Total/total*100)}%`}/>
                          </Box>
                        </Stack>
                        <LinearProgress sx={{backgroundColor: '#bfbfbf', '& .MuiLinearProgress-bar': {backgroundColor: '#cf473bf5'}}}
                          color="primary" variant="determinate" value={Math.round(value.Total/total*100)}/>
                      </Box>
                    </ListItem>
                  )
                }
                </React.Fragment>
              }
            </List>
          </CardContent>
        </Grid>
      </Grid>
    </Card>
  );
}

