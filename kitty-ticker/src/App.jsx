import { useState, useEffect } from 'react'
import styled, { keyframes } from 'styled-components';
import CurrencyCard from './components/CurrencyCard';
import PartyConfetti from './components/PartyConfetti';
import { getCurrencyPairs } from './services/tickerService';
import './App.css'

// Party backgrounds and animations
const partyColors = keyframes`
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
`;

const floatAnimation = keyframes`
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
  100% {
    transform: translateY(0px);
  }
`;

const AppContainer = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: linear-gradient(270deg, #ff9a9e, #fad0c4, #fad0c4, #a1c4fd, #c2e9fb);
  background-size: 1000% 1000%;
  animation: ${partyColors} 30s ease infinite;
  padding: 2rem;
`;

const Header = styled.header`
  text-align: center;
  margin-bottom: 2rem;
`;

const Title = styled.h1`
  font-family: 'Pacifico', cursive;
  font-size: 3.5rem;
  color: white;
  margin-bottom: 0.5rem;
  text-shadow: 3px 3px 0 #ff9a9e,
               6px 6px 0 rgba(0,0,0,0.2);
  animation: ${floatAnimation} 3s ease-in-out infinite;
`;

const Subtitle = styled.p`
  font-size: 1.2rem;
  color: white;
  margin: 0;
  text-shadow: 1px 1px 4px rgba(0,0,0,0.2);
`;

const CardContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
`;

const Footer = styled.footer`
  margin-top: 3rem;
  color: white;
  font-size: 0.9rem;
  text-align: center;
  width: 100%;
`;

function App() {
  const [currencyPairs, setCurrencyPairs] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Get the currency pairs when the component mounts
    setCurrencyPairs(getCurrencyPairs());
    setIsLoading(false);
  }, []);

  return (
    <AppContainer>
      <PartyConfetti count={30} />
      <Header>
        <Title>Model Context Party</Title>
        <Subtitle>Live Cryptocurrency Tickers with Cats</Subtitle>
      </Header>

      {isLoading ? (
        <p>Loading currency data...</p>
      ) : (
        <CardContainer>
          {currencyPairs.map((pair) => (
            <CurrencyCard key={pair} pair={pair} refreshInterval={10000} />
          ))}
        </CardContainer>
      )}

      <Footer>
        <p>Created for Model Context Party Hackathon | {new Date().getFullYear()}</p>
        <p>Powered by Luno API</p>
      </Footer>
    </AppContainer>
  )
}

export default App
