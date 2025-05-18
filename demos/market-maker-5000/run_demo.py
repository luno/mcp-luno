#!/usr/bin/env python3
"""
Demo script to showcase the Luno Market Maker Bot functionality.
This script places a single market making order on the XBTZAR pair
and displays the result.
"""

import os
import json
import time
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Import functions from the main bot
from market_maker_bot import (
    log, get_balances, get_ticker, get_order_book,
    create_order, get_market_making_prices
)

def run_demo():
    """Run a demonstration of the market maker bot by placing a single order."""
    print("=" * 80)
    print("LUNO MARKET MAKER BOT - DEMO")
    print("=" * 80)

    # Select a trading pair
    pair = "XBTZAR"

    print(f"\nRunning demo for {pair}...\n")

    # Get account balances
    print("Fetching account balances...")
    balances = get_balances()
    if not balances:
        print("Failed to get account balances. Please check your API credentials.")
        return

    # Print relevant balances
    base_asset = "XBT"  # For XBTZAR
    quote_asset = "ZAR"  # For XBTZAR

    base_balance = next((b for b in balances if b["asset"] == base_asset), None)
    quote_balance = next((b for b in balances if b["asset"] == quote_asset), None)

    if base_balance and quote_balance:
        print(f"{base_asset} Balance: {float(base_balance.get('balance', 0))}")
        print(f"{quote_asset} Balance: {float(quote_balance.get('balance', 0))}")

    # Get market data
    print("\nFetching market data...")
    ticker = get_ticker(pair)
    if not ticker:
        print("Failed to get market data. Please check the trading pair.")
        return

    print(f"Last trade price: {ticker.get('last_trade')}")
    print(f"Current bid: {ticker.get('bid')}")
    print(f"Current ask: {ticker.get('ask')}")

    # Calculate market making prices
    prices = get_market_making_prices(ticker)
    print("\nCalculated market making prices:")
    print(f"Bid price: {prices['bid']}")
    print(f"Ask price: {prices['ask']}")

    # Get order book
    order_book = get_order_book(pair)
    print("\nCurrent order book (top 3 entries):")
    if order_book:
        print("Asks (Sell orders):")
        for i, ask in enumerate(order_book.get("asks", [])[:3]):
            print(f"  {ask.get('volume')} @ {ask.get('price')}")

        print("Bids (Buy orders):")
        for i, bid in enumerate(order_book.get("bids", [])[:3]):
            print(f"  {bid.get('volume')} @ {bid.get('price')}")

    # Place a small demo order
    print("\nPlacing a demo order...")

    # Calculate a small order size
    small_volume = 0.001  # A very small amount of XBT

    # Check if we have enough balance
    if base_balance and float(base_balance.get("balance", 0)) > small_volume:
        print(f"Placing a sell order for {small_volume} {base_asset} @ {prices['ask']} {quote_asset}")

        # Create the order
        order_id = create_order(pair, "ASK", small_volume, prices["ask"])

        if order_id:
            print(f"\nOrder successfully placed!")
            print(f"Order ID: {order_id}")
            print("\nDemo completed successfully!")
        else:
            print("\nFailed to place the order. Check the error logs.")
    else:
        print(f"Not enough {base_asset} balance to place a demo order.")
        print("Demo completed without placing an order.")

if __name__ == "__main__":
    run_demo()
