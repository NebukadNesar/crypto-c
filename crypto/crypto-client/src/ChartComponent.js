import React, {useEffect, useState} from 'react';
import {Line} from 'react-chartjs-2';
import {CategoryScale, Chart as ChartJS, Legend, LinearScale, LineElement, PointElement, Tooltip} from 'chart.js';

// Register necessary components
ChartJS.register(LineElement, CategoryScale, LinearScale, PointElement, Legend, Tooltip);

const ChartComponent = ({cryptoData}) => {
    const [chartData, setChartData] = useState(null);

    useEffect(() => {
        if (cryptoData && cryptoData.length > 0) {
            // Map timestamps and prices
            const labels = cryptoData.map(item => new Date(item.Timestamp).toLocaleTimeString());
            const prices = cryptoData.map(item => parseFloat(item.Data.priceUsd));

            // Set chart data
            setChartData({
                labels: labels, // X-axis: Timestamps
                datasets: [{
                    label: 'Bitcoin Price (USD)', data: prices, borderColor: 'red', fill: false,
                }],
            });
        }
    }, [cryptoData]); // Runs when cryptoData changes

    const options = {
        responsive: true, plugins: {
            legend: {
                display: true,
            },
        },
    };

    return (<div style={{width: '100%', maxWidth: '100%', height: '300px'}}>
            {chartData ? <Line data={chartData} options={options}/> : <p>Loading data...</p>}
        </div>);
};

export default ChartComponent;
