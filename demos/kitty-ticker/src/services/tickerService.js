// This is a mock service since we're using sample data from #get_ticker
// In a real application, you would connect to the Luno API

const API_ENDPOINT = 'https://api.luno.com/api/1';

// The 5 most popular cryptocurrency pairs in South Africa
const CURRENCY_PAIRS = [
  'XBTZAR', // Bitcoin to South African Rand
  'ETHZAR', // Ethereum to South African Rand
  'LTCZAR', // Litecoin to South African Rand
  'XRPZAR', // Ripple to South African Rand
  'SOLZAR'  // Solana to South African Rand
];

// Sample data based on the #get_ticker attachment
const SAMPLE_DATA = {
  'XBTZAR': {
    ask: "1892075.00",
    bid: "1885413.00",
    last_trade: "1888574.00",
    pair: "XBTZAR",
    rolling_24_hour_volume: "1.370151",
    status: "ACTIVE",
    timestamp: 1747499349235
  },
  'ETHZAR': {
    ask: "105207.00",
    bid: "104932.00",
    last_trade: "105000.00",
    pair: "ETHZAR",
    rolling_24_hour_volume: "15.23456",
    status: "ACTIVE",
    timestamp: 1747499349235
  },
  'LTCZAR': {
    ask: "3125.50",
    bid: "3120.75",
    last_trade: "3123.00",
    pair: "LTCZAR",
    rolling_24_hour_volume: "42.12345",
    status: "ACTIVE",
    timestamp: 1747499349235
  },
  'XRPZAR': {
    ask: "21.75",
    bid: "21.65",
    last_trade: "21.70",
    pair: "XRPZAR",
    rolling_24_hour_volume: "3500.67890",
    status: "ACTIVE",
    timestamp: 1747499349235
  },
  'SOLZAR': {
    ask: "3450.25",
    bid: "3445.50",
    last_trade: "3448.00",
    pair: "SOLZAR",
    rolling_24_hour_volume: "78.45678",
    status: "ACTIVE",
    timestamp: 1747499349235
  }
};

export const getTicker = async (pair) => {
  try {
    // In a real application, this would be:
    // const response = await axios.get(`${API_ENDPOINT}/ticker?pair=${pair}`);
    // return response.data;

    // For now, return our sample data
    return SAMPLE_DATA[pair];
  } catch (error) {
    console.error(`Error fetching ticker for ${pair}:`, error);
    throw error;
  }
};

export const getAllTickers = async () => {
  try {
    // In a real application, we would fetch all tickers from the API
    // For now, return our sample data for all pairs
    const tickers = {};
    for (const pair of CURRENCY_PAIRS) {
      tickers[pair] = SAMPLE_DATA[pair];
    }
    return tickers;
  } catch (error) {
    console.error('Error fetching all tickers:', error);
    throw error;
  }
};

export const getCurrencyPairs = () => {
  return CURRENCY_PAIRS;
};
