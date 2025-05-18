# Luno Market Maker 5000

```
 __       __                      __                   __
|  \     /  \                    |  \                 |  \
| $$\   /  $$  ______    ______  | $$   __   ______  _| $$_
| $$$\ /  $$$ |      \  /      \ | $$  /  \ /      \|   $$ \
| $$$$\  $$$$  \$$$$$$\|  $$$$$$\| $$_/  $$|  $$$$$$\\$$$$$$
| $$\$$ $$ $$ /      $$| $$   \$$| $$   $$ | $$    $$ | $$ __
| $$ \$$$| $$|  $$$$$$$| $$      | $$$$$$\ | $$$$$$$$ | $$|  \
| $$  \$ | $$ \$$    $$| $$      | $$  \$$\ \$$     \  \$$  $$
 \$$      \$$  \$$$$$$$ \$$       \$$   \$$  \$$$$$$$   \$$$$

 __       __            __
|  \     /  \          |  \
| $$\   /  $$  ______  | $$   __   ______    ______
| $$$\ /  $$$ |      \ | $$  /  \ /      \  /      \
| $$$$\  $$$$  \$$$$$$\| $$_/  $$|  $$$$$$\|  $$$$$$\
| $$\$$ $$ $$ /      $$| $$   $$ | $$    $$| $$   \$$
| $$ \$$$| $$|  $$$$$$$| $$$$$$\ | $$$$$$$$| $$
| $$  \$ | $$ \$$    $$| $$  \$$\ \$$     \| $$
 \$$      \$$  \$$$$$$$ \$$   \$$  \$$$$$$$ \$$

 _______    ______    ______    ______
|       \  /      \  /      \  /      \
| $$$$$$$ |  $$$$$$\|  $$$$$$\|  $$$$$$\
| $$____  | $$$\| $$| $$$\| $$| $$$\| $$
| $$    \ | $$$$\ $$| $$$$\ $$| $$$$\ $$
 \$$$$$$$\| $$\$$\$$| $$\$$\$$| $$\$$\$$
|  \__| $$| $$_\$$$$| $$_\$$$$| $$_\$$$$
 \$$    $$ \$$  \$$$ \$$  \$$$ \$$  \$$$
  \$$$$$$   \$$$$$$   \$$$$$$   \$$$$$$
```

This is an automated market maker bot for the Luno cryptocurrency exchange. The bot places buy and sell orders around the market price to provide liquidity and potentially profit from the bid-ask spread.

## Features

- Trades the top 5 cryptocurrency pairs on Luno based on available balances
- Implements a basic market making strategy with configurable parameters
- Automatically refreshes orders at regular intervals
- Includes safety measures to prevent excessive trading or losses

## Strategy

The bot uses a simple market making strategy:

1. For each trading pair, it fetches the current market price
2. It calculates a bid price slightly below and an ask price slightly above the mid-market price
3. Places limit orders on both sides of the order book
4. Refreshes orders periodically to adjust to market movements
5. Uses a percentage of the available balance for each order to manage risk

## Setup

1. Install the required dependencies:
   ```
   pip install -r requirements.txt
   ```

2. Create a `.env` file in the root directory with your Luno API credentials:
   ```
   LUNO_API_KEY=your_api_key_here
   LUNO_API_SECRET=your_api_secret_here
   ```

3. Configure the bot parameters in `market_maker_bot.py` (optional):
   - `TRADING_PAIRS`: List of pairs to trade
   - `SPREAD_PERCENTAGE`: How wide of a spread to maintain around mid price
   - `ORDER_SIZE_PERCENTAGE`: Percentage of available balance to use per order
   - `ORDER_EXPIRY_MINUTES`: How frequently to refresh orders

## Usage

Run the bot with:

```
python market_maker_bot.py
```

The bot will continuously run until you stop it with Ctrl+C. All orders will be canceled when the bot is stopped.
