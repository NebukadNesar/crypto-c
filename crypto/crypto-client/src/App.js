import './App.css';
import {useEffect, useRef, useState} from "react";
import CryptoDetail from "./CryptoDetail";
import ChartComponent from './ChartComponent'

function App() {
    const ws = useRef(null);

    useEffect(() => {
        // Create WebSocket connection.
        ws.current = new WebSocket("ws://localhost:9988/crypto");

        // Connection opened
        ws.current.addEventListener("open", (event) => {
            console.log("Connected to the server ", event);
        });

        // Listen for messages
        ws.current.addEventListener("message", (event) => {
            let message = JSON.parse(event.data);
            if (message.id === 'snapshot') {
                setCryptoData(message.data)
            } else if (message.id === 'historical') {
                setChartsData(message.data)
            } else {
                console.log("no such message supported", message)
            }
        });
    }, []);

    const [cryptoData, setCryptoData] = useState(null);
    const [chartData, setChartsData] = useState(null);

    return (<div className="App">
        <h1>Crypto Data</h1>
        <ChartComponent cryptoData={chartData}/>
        <div>
            <h2>Coins</h2>
            {cryptoData && cryptoData.length > 0 ? (<div style={styles.cryptoList}>
                {cryptoData.map((crypto) => (<CryptoDetail key={crypto.id} cryptoData={crypto}/>))}
            </div>) : (<p>No data available</p>)}
        </div>
    </div>);
}
export default App;
