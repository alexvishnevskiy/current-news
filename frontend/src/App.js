import React from 'react';
import axios from 'axios'

class App extends React.Component{
  constructor(props){
    super(props);
    this.state = {
      topics: []
    }
  }

  componentDidMount() {
    axios.get('http://localhost:8888/data', {mode: 'no-cors'})
    .then(res => this.setState({topics: res.data.Europe}));
  }

  render() {//returns [] the first render, returns ['topic1','topic2','topic3'] on the second render;
    return(
        <ul>
          {this.state.topics.map(topic => (
              <li key = {topic}>{topic}</li>
          ))}
        </ul>
    )
  }
}

export default App;
