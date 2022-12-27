import React, { useState, useEffect} from 'react';
import axios from 'axios';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import Chip from '@mui/material/Chip';
import Card from '@mui/material/Card';
import Paper from '@mui/material/Paper';
import InputBase from '@mui/material/InputBase';
import SearchIcon from '@mui/icons-material/Search';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CardContent from '@mui/material/CardContent';
import ArticleIcon from '@mui/icons-material/Article';
import { FixedSizeList } from 'react-window';

const Row = props => {
  const { data, index, style } = props;

  if (data.testData[index].Source != "") {
    return (
      <ListItem style={style} key={index} component="div" disablePadding>
          <ListItemText primary={data.testData[index].Title} />
            <Chip label={data.testData[index].Source} color="error" />
          <IconButton type="button" href={data.testData[index].URL}>
              <ArticleIcon/>
          </IconButton>
      </ListItem>
    );
  } 
  return (
    <ListItem style={style} key={index} component="div" disablePadding>
          <ListItemText primary={data.testData[index].Title} />
          <IconButton type="button" href={data.testData[index].URL}>
              <ArticleIcon/>
          </IconButton>
      </ListItem>
  )
}

export default function Headlines() {
  const [dict, updataData] = useState({
    question: "",
    articles: {},
  });
   
  useEffect(() => {
    axios.get('http://localhost:8080/headlines', {mode: 'no-cors'})
    // axios.get('http://4.227.248.158/data', {mode: 'no-cors'})
    .then(res => res.data)
    .then(res => {
      updataData({
        question: "",
        articles: res
      });
    });
    }, []);

  const onChange = (event) => {
    updataData({
      question: event.target.value,
      articles: dict.articles
    });
  };

  const handleClick = () => {
    axios.get(`http://localhost:8080/headlines/${dict.question}`, {mode: 'no-cors'})
    .then(res => res.data)
    .then(res => {
      updataData({
        question: "",
        articles: res
      });
    });
  };

  return (
    <Card sx={{
      width: '100%',
      minHeight: 300,
    }}>
        <CardContent>
            <Typography align="left" variant="h8" component="div">
                <b>Top headlines</b>
            </Typography>
        </CardContent>
        <Paper component="form"
        sx={{ p: '2px 4px', display: 'flex', alignItems: 'center'}}>
            <InputBase
                sx={{ ml: 1, flex: 1}}
                placeholder="Search article topic"
                onChange={onChange}
                inputProps={{ 'aria-label': 'search article topic' }}
            />
            <IconButton type="button" onClick={handleClick} sx={{ p: '10px' }} aria-label="search">
                <SearchIcon />
            </IconButton>
        </Paper>
        <FixedSizeList
            height={267}
            // width={450}
            itemSize={80}
            itemCount={dict.articles && dict.articles.length}
            overscanCount={5}
            itemData={{ testData: dict.articles}}
        >
            {Row}
        </FixedSizeList>
    </Card>
  );
}