import React, { useEffect, useState } from 'react';
import styled, { keyframes } from 'styled-components';
import { getTicker } from '../services/tickerService';
import PropTypes from 'prop-types';

// Import the correct images from src/assets/images
import kitty1 from '../assets/images/5a521abe2f93c7a8d5137fa1.png';
import kitty2 from '../assets/images/5a521ac92f93c7a8d5137fa3.png';
import kitty3 from '../assets/images/5a521ad42f93c7a8d5137fa5.png';
import kitty4 from '../assets/images/5a521adb2f93c7a8d5137fa6.png';
import kitty5 from '../assets/images/5a521ae02f93c7a8d5137fa7.png';

const catImages = {
  'XBTZAR': kitty1,
  'ETHZAR': kitty2,
  'LTCZAR': kitty3,
  'XRPZAR': kitty4,
  'SOLZAR': kitty5
};

const currencyNames = {
  'XBTZAR': 'Bitcoin',
  'ETHZAR': 'Ethereum',
  'LTCZAR': 'Litecoin',
  'XRPZAR': 'Ripple',
  'SOLZAR': 'Solana'
};

const pulse = keyframes`
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
`;

const Card = styled.div`
  background: ${props => {
    if (props.theme === 'up') return 'linear-gradient(to bottom, #a8e063, #56ab2f)';
    if (props.theme === 'down') return 'linear-gradient(to bottom, #ff9966, #ff5e62)';
    return 'linear-gradient(to bottom, #4e54c8, #8f94fb)';
  }};
  border-radius: 16px;
  padding: 20px;
  margin: 15px;
  box-shadow: 0 10px 20px rgba(0,0,0,0.19), 0 6px 6px rgba(0,0,0,0.23);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  width: 300px;
  height: 400px;
  position: relative;
  overflow: hidden;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 30px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
  }
`;

const CardContent = styled.div`
  z-index: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
`;

const CatImage = styled.img`
  width: 100%;
  height: 180px;
  object-fit: cover;
  border-radius: 12px;
  margin-bottom: 15px;
  transition: transform 0.3s ease;

  &:hover {
    transform: scale(1.05);
  }
`;

const CurrencyName = styled.h2`
  color: white;
  margin: 0;
  font-size: 24px;
  text-shadow: 1px 1px 3px rgba(0,0,0,0.3);
`;

const CurrencyPair = styled.h3`
  color: rgba(255,255,255,0.8);
  margin: 5px 0 15px 0;
  font-size: 16px;
`;

const PriceContainer = styled.div`
  margin-top: auto;
`;

const Price = styled.div`
  font-size: 26px;
  font-weight: bold;
  color: white;
  margin-bottom: 5px;
  animation: ${props => props.priceChanged ? pulse : 'none'} 0.5s;
`;

const PriceChange = styled.div`
  display: flex;
  align-items: center;
  margin-top: 8px;
  color: white;
`;

const Volume = styled.div`
  font-size: 14px;
  color: rgba(255,255,255,0.8);
  margin-top: 5px;
`;

const LastUpdated = styled.div`
  font-size: 12px;
  color: rgba(255,255,255,0.7);
  margin-top: 10px;
`;

const CurrencyCard = ({ pair, refreshInterval = 30000 }) => {
  const [tickerData, setTickerData] = useState(null);
  const [priceChanged, setPriceChanged] = useState(false);
  const [theme, setTheme] = useState('neutral');
  const [previousPrice, setPreviousPrice] = useState(null);

  useEffect(() => {
    // Fetch initial data
    const fetchTickerData = async () => {
      try {
        const data = await getTicker(pair);
        if (tickerData && tickerData.last_trade !== data.last_trade) {
          setPriceChanged(true);
          setTimeout(() => setPriceChanged(false), 1000);
          if (parseFloat(data.last_trade) > parseFloat(tickerData.last_trade)) {
            setTheme('up');
          } else if (parseFloat(data.last_trade) < parseFloat(tickerData.last_trade)) {
            setTheme('down');
          }
          setPreviousPrice(tickerData.last_trade);
        }
        setTickerData(data);
      } catch (error) {
        console.error(`Error fetching data for ${pair}:`, error);
      }
    };
    fetchTickerData();
    const intervalId = setInterval(fetchTickerData, refreshInterval);
    return () => clearInterval(intervalId);
  }, [pair, refreshInterval, tickerData]);

  if (!tickerData) {
    return (
      <Card theme="neutral">
        <CardContent>
          <CurrencyName>Loading...</CurrencyName>
        </CardContent>
      </Card>
    );
  }

  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-ZA', {
      style: 'currency',
      currency: 'ZAR',
      minimumFractionDigits: 2
    }).format(price);
  };

  const formatDateTime = (timestamp) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  console.log('catImages:', catImages);
  console.log('pair:', pair);
  console.log('catImages[pair]:', catImages[pair]);

  return (
    <Card theme={theme}>
      <CardContent>
        <CatImage src={catImages[pair]} alt={`${currencyNames[pair]} cat`} onError={e => {console.error('Image error', e);}} />
        <CurrencyName>{currencyNames[pair]}</CurrencyName>
        <CurrencyPair>{pair}</CurrencyPair>

        <PriceContainer>
          <Price priceChanged={priceChanged}>
            {formatPrice(tickerData.last_trade)}
          </Price>

          {previousPrice && (
            <PriceChange>
              {parseFloat(tickerData.last_trade) > parseFloat(previousPrice) ? '↑' : '↓'}
              {' '}
              {formatPrice(Math.abs(parseFloat(tickerData.last_trade) - parseFloat(previousPrice)))}
            </PriceChange>
          )}

          <Volume>
            Volume: {parseFloat(tickerData.rolling_24_hour_volume).toFixed(4)}
          </Volume>

          <LastUpdated>
            Last updated: {formatDateTime(tickerData.timestamp)}
          </LastUpdated>
        </PriceContainer>
      </CardContent>
    </Card>
  );
};

CurrencyCard.propTypes = {
  pair: PropTypes.string.isRequired,
  refreshInterval: PropTypes.number
};

export default CurrencyCard;
