import Chart from "react-apexcharts";
import Card from '@mui/material/Card';
import axios from 'axios';
import Typography from '@mui/material/Typography';
import CardContent from '@mui/material/CardContent';
import React, { useState, useEffect} from 'react';


function processData(data) {
    var dict = {
        'timestamps': [],
        'values': []
    };
    for (var i = 0; i < data.archive.length; i++) {
        dict['timestamps'].push(data.archive[i].Time)
        dict['values'].push(data.archive[i].Value)
    }
    return dict
}

export default function ArchiveChart() {    
    const [data, updataData] = useState({
        options: {
            stroke: {
                curve: 'smooth',
            },
            xaxis: {
                categories: []
            },
            colors: ['#cf473bf5']
        },
        series: [
        {
        name: "Number of articles",
        data: []
        }
        ]
    });

    useEffect(() => {
        axios.get('http://localhost:8080/archive', {mode: 'no-cors'})
        .then(res => res.data)
        .then(res => processData(res))
        .then(res => {
            updataData({
                options: {
                    stroke: {
                        curve: 'smooth',
                    },
                    xaxis: {
                        categories: res.timestamps
                    },
                    colors: ['#cf473bf5']
                },
                series: [
                    {
                    name: "Number of articles",
                    data: res.values
                    }
                ]
            });
        })
        }, []);

    return (
        <Card>
            <CardContent>
                <Typography align="center" variant="h7" component="div">
                    <b>Number of NY Times articles</b>
                </Typography>
             </CardContent>
            <Chart
                    options={data.options}
                    series={data.series}
                    type="line"
                    height="300"
                />
        </Card>
    );
}