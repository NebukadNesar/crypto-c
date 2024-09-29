import React from 'react';

const CryptoDetail = ({ cryptoData }) => {
    return (
        //todo format all fields and make it list view instead
        <div style={styles.cryptoDetail}>
            <h1>{cryptoData.name} ({cryptoData.symbol})</h1>
            <p>Rank: {cryptoData.rank}</p>
            <p>Supply: {cryptoData.supply}</p>
            <p>Max Supply: {cryptoData.maxSupply}</p>
            <p>Market Cap (USD): {cryptoData.marketCapUsd}</p>
            <p>24h Volume (USD): {cryptoData.volumeUsd24Hr}</p>
            <p>Current Price (USD): {cryptoData.priceUsd}</p>
            <p>Change (24h): {cryptoData.changePercent24Hr}%</p>
            <p>VWAP (24h): {cryptoData.vwap24Hr}</p>
            <div style={styles.explorerContainer}>
                Explorer:
                <a
                    href={cryptoData.explorer}
                    target="_blank"
                    rel="noopener noreferrer"
                    title="View on Explorer"
                    style={styles.explorerLink}
                >
                    Visit Explorer
                </a>
            </div>
        </div>
    );
};

// Styles for the CryptoDetail component
const styles = {
    cryptoDetail: {
        padding: '20px',
        margin: '10px',
        backgroundColor: '#fff',
        borderRadius: '5px',
        boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
        flex: '1 1 calc(20% - 20px)', // Adjust as needed for layout
    },
    explorerContainer: {
        marginTop: '10px',
    },
    explorerLink: {
        color: '#007bff',
        textDecoration: 'none',
    },
};

export default CryptoDetail;
