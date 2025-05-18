#!/usr/bin/env python3
"""
Test module for the Luno Market Maker Bot.
This script allows you to test various components of the market maker bot
without placing actual orders.
"""

import argparse
import json
import os
from datetime import datetime
from dotenv import load_dotenv

# Import functions from the main bot
from market_maker_bot import (
    log, get_balances, get_ticker, get_order_book,
    get_open_orders, get_market_making_prices,
    calculate_order_size, get_asset_balance
)

# Load environment variables from .env file
load_dotenv()

def print_header(title):
    """Print a header for the output."""
    print("\n" + "=" * 80)
    print(f" {title}")
    print("=" * 80)

def print_json(data):
    """Print JSON data in a pretty format."""
    if data:
        print(json.dumps(data, indent=2))
    else:
        print("No data returned")

def test_balance_check():
    """Test retrieving and displaying account balances."""
    print_header("ACCOUNT BALANCES")

    balances = get_balances()
    if balances:
        # Sort by balance value descending
        sorted_balances = sorted(
            balances,
            key=lambda x: float(x.get("balance", 0)),
            reverse=True
        )

        print("Asset".ljust(10), "Balance".ljust(20), "Reserved".ljust(20), "Available".ljust(20))
        print("-" * 70)

        for account in sorted_balances:
            asset = account.get("asset", "")
            balance = float(account.get("balance", 0))
            reserved = float(account.get("reserved", 0))
            available = balance - reserved

            # Only show assets with non-zero balance
            if balance > 0:
                print(
                    f"{asset}".ljust(10),
                    f"{balance:.8f}".ljust(20),
                    f"{reserved:.8f}".ljust(20),
                    f"{available:.8f}".ljust(20)
                )

def test_market_data(pair):
    """Test retrieving market data for a specific pair."""
    print_header(f"MARKET DATA FOR {pair}")

    # Get ticker data
    print("\nTicker:")
    ticker = get_ticker(pair)
    if ticker:
        print(f"Last trade: {ticker.get('last_trade', 'N/A')}")
        print(f"Bid: {ticker.get('bid', 'N/A')}")
        print(f"Ask: {ticker.get('ask', 'N/A')}")
        print(f"24h volume: {ticker.get('rolling_24_hour_volume', 'N/A')}")

        # Calculate market making prices
        prices = get_market_making_prices(ticker)
        print("\nMarket Making Prices:")
        print(f"Bid (Buy) price: {prices['bid']}")
        print(f"Ask (Sell) price: {prices['ask']}")

    # Get order book data
    print("\nOrder Book (Top 5 entries):")
    order_book = get_order_book(pair)
    if order_book:
        print("Asks (Sell orders):")
        for i, ask in enumerate(order_book.get("asks", [])[:5]):
            print(f"  {ask.get('volume')} @ {ask.get('price')}")

        print("Bids (Buy orders):")
        for i, bid in enumerate(order_book.get("bids", [])[:5]):
            print(f"  {bid.get('volume')} @ {bid.get('price')}")

def test_order_sizing(pair):
    """Test order size calculations for a specific pair."""
    print_header(f"ORDER SIZE CALCULATION FOR {pair}")

    balances = get_balances()
    ticker = get_ticker(pair)

    if balances and ticker:
        bid_price = float(ticker.get("bid", 0))
        ask_price = float(ticker.get("ask", 0))

        # Calculate mid price
        mid_price = (bid_price + ask_price) / 2 if bid_price and ask_price else 0

        if mid_price > 0:
            # Calculate buy order size
            buy_size = calculate_order_size(pair, balances, mid_price, "BUY")
            print(f"Buy order size at price {mid_price}: {buy_size:.8f}")

            # Calculate sell order size
            sell_size = calculate_order_size(pair, balances, mid_price, "ASK")
            print(f"Sell order size at price {mid_price}: {sell_size:.8f}")
        else:
            print("Could not calculate order sizes, invalid price data")

def test_open_orders(pair=None):
    """Test retrieving open orders."""
    print_header("OPEN ORDERS" + (f" FOR {pair}" if pair else ""))

    orders = get_open_orders(pair)
    if orders and "orders" in orders:
        orders_list = orders["orders"]
        if orders_list:
            for order in orders_list:
                print(f"Order ID: {order.get('order_id')}")
                print(f"  Pair: {order.get('pair', 'N/A')}")
                print(f"  Type: {order.get('type', 'N/A')}")
                print(f"  Price: {order.get('price', 'N/A')}")
                print(f"  Volume: {order.get('volume', 'N/A')}")
                print(f"  Created: {order.get('creation_timestamp', 'N/A')}")
                print()
        else:
            print("No open orders found.")
    else:
        print("Failed to retrieve open orders.")

def main():
    """Main function to run tests."""
    parser = argparse.ArgumentParser(description="Test the Luno Market Maker Bot")
    parser.add_argument("--balances", action="store_true", help="Test account balance retrieval")
    parser.add_argument("--market", help="Test market data for a specific pair (e.g., XBTZAR)")
    parser.add_argument("--orders", help="Test order sizing for a specific pair (e.g., XBTZAR)")
    parser.add_argument("--open-orders", action="store_true", help="List all open orders")
    parser.add_argument("--all", action="store_true", help="Run all tests")

    args = parser.parse_args()

    # If no arguments were provided, show help and exit
    if not any(vars(args).values()):
        parser.print_help()
        return

    # Execute requested tests
    if args.balances or args.all:
        test_balance_check()

    if args.market or args.all:
        pairs = [args.market] if args.market else ["XBTZAR", "ETHZAR", "XRPZAR", "ETHXBT", "USDCZAR"]
        for pair in pairs:
            test_market_data(pair)

    if args.orders or args.all:
        pairs = [args.orders] if args.orders else ["XBTZAR", "ETHZAR", "XRPZAR", "ETHXBT", "USDCZAR"]
        for pair in pairs:
            test_order_sizing(pair)

    if args.open_orders or args.all:
        test_open_orders(args.open_orders if isinstance(args.open_orders, str) else None)

if __name__ == "__main__":
    main()
