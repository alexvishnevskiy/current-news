import * as React from 'react';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import Card from '@mui/material/Card';
import Paper from '@mui/material/Paper';
import InputBase from '@mui/material/InputBase';
import SearchIcon from '@mui/icons-material/Search';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CardContent from '@mui/material/CardContent';
import ArticleIcon from '@mui/icons-material/Article';
import Divider from '@mui/material/Divider';
import { FixedSizeList } from 'react-window';

const Row = props => {
  const { data, index, style } = props;

  return (
    <ListItem style={style} key={index} component="div" disablePadding>
        <ListItemText primary={data.testData[index].title} />
        <IconButton type="button" href={data.testData[index].url}>
            <ArticleIcon/>
        </IconButton>
    </ListItem>
  );
}

export default function Headlines() {
//   const [articles, updataData] = useState(null);
  const articles = [{"title": "I love you and I need you, vdfvdvdf vfdvdfvdf vf vdf dfvdf vdfv vfd vfd vdf vfd vfd vf vfd vf I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"},
  {"title": "I love you and I need you, I love you and I need you. I love you and I need you, I love you and I need you", "url": "https://www.robinwieruch.de/react-pass-props-to-component/"}
]
   
  // useEffect(() => {
  //   axios.get('http://localhost:8080/headlines', {mode: 'no-cors'})
  //   // axios.get('http://4.227.248.158/data', {mode: 'no-cors'})
  //   .then(res => res.data)
  //   .then(res => {
  //       updataData(res);
  //     });
  //   }, []);

  return (
    <Card>
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
                inputProps={{ 'aria-label': 'search article topic' }}
            />
            <IconButton type="button" sx={{ p: '10px' }} aria-label="search">
                <SearchIcon />
            </IconButton>
        </Paper>
        <FixedSizeList
            height={300}
            // width={450}
            itemSize={120}
            itemCount={articles.length}
            overscanCount={5}
            itemData={{ testData: articles}}
        >
            {Row}
        </FixedSizeList>
    </Card>
  );
}