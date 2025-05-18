#!/usr/bin/env python3
"""
Luno Market Maker Bot
---------------------
This bot implements a basic market making strategy for the top cryptocurrency pairs on Luno.
It creates orders with small spreads around the current market price to provide liquidity
and potentially profit from the bid-ask spread.
"""

import requests
import json
import time
import os
import random
from datetime import datetime
import hmac
import hashlib

# Configuration
API_KEY = os.environ.get("LUNO_API_KEY")
API_SECRET = os.environ.get("LUNO_API_SECRET")
BASE_URL = "https://api.luno.com/api/1"

# Trading pairs to focus on (top 5 based on available balances)
TRADING_PAIRS = [
    "XBTZAR", "ETHZAR", "XRPZAR", "ETHXBT", "USDCZAR"
]

# Market making parameters
SPREAD_PERCENTAGE = 0.5  # 0.5% spread around mid price
ORDER_SIZE_PERCENTAGE = 2.0  # Use 2% of available balance for each order
ORDER_EXPIRY_MINUTES = 5  # Cancel and replace orders every 5 minutes
MAX_ORDERS_PER_PAIR = 2  # Maximum number of orders (1 buy, 1 sell) per pair
PRICE_PRECISION = 0  # Decimal places for price
VOLUME_PRECISION = 6  # Decimal places for volume

# Logging configuration
def log(message):
    """Simple logging function that prints messages with timestamps."""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"[{timestamp}] {message}")

def make_request(method, endpoint, params=None, data=None):
    """Make an authenticated request to the Luno API."""
    url = f"{BASE_URL}/{endpoint}"
    headers = {"Content-Type": "application/x-www-form-urlencoded"}

    if method == "GET":
        response = requests.get(url, params=params, auth=(API_KEY, API_SECRET), headers=headers)
    else:
        response = requests.post(url, data=data, auth=(API_KEY, API_SECRET), headers=headers)

    if response.status_code not in [200, 201, 202]:
        log(f"API Error: {response.status_code} - {response.text}")
        return None

    return response.json()

def get_balances():
    """Get account balances."""
    return make_request("GET", "accounts")

def get_ticker(pair):
    """Get ticker information for a trading pair."""
    return make_request("GET", "ticker", params={"pair": pair})

def get_order_book(pair):
    """Get order book for a trading pair."""
    return make_request("GET", "orderbook", params={"pair": pair})

def create_order(pair, order_type, volume, price):
    """Create a limit order."""
    data = {
        "pair": pair,
        "type": order_type,
        "volume": f"{volume:.{VOLUME_PRECISION}f}",
        "price": f"{price:.{PRICE_PRECISION}f}",
        "post_only": "true"  # Ensure orders are only posted to the order book, never taken
    }

    response = make_request("POST", "postorder", data=data)
    if response and "order_id" in response:
        log(f"Created {order_type} order for {pair}: {volume:.{VOLUME_PRECISION}f} @ {price:.{PRICE_PRECISION}f} (ID: {response['order_id']})")
        return response["order_id"]
    else:
        log(f"Failed to create {order_type} order for {pair}")
        return None

def cancel_order(order_id):
    """Cancel an existing order."""
    data = {"order_id": order_id}
    response = make_request("POST", "stoporder", data=data)
    if response and response.get("success", False):
        log(f"Cancelled order: {order_id}")
        return True
    else:
        log(f"Failed to cancel order: {order_id}")
        return False

def get_open_orders(pair=None):
    """Get all open orders, optionally filtered by pair."""
    params = {}
    if pair:
        params["pair"] = pair

    return make_request("GET", "listorders", params=params)

def cancel_all_orders(pair=None):
    """Cancel all open orders, optionally filtered by pair."""
    orders = get_open_orders(pair)
    if not orders or "orders" not in orders:
        return

    for order in orders["orders"]:
        cancel_order(order["order_id"])

def get_asset_balance(balances, asset):
    """Extract balance for a specific asset."""
    for account in balances:
        if account["asset"] == asset:
            return {
                "balance": float(account["balance"]),
                "reserved": float(account["reserved"]),
                "available": float(account["balance"]) - float(account["reserved"])
            }
    return {"balance": 0, "reserved": 0, "available": 0}

def calculate_order_size(pair, balances, price, side):
    """Calculate the order size based on available balance."""
    if side == "BUY":
        quote_asset = pair[3:]  # e.g., ZAR in XBTZAR
        quote_balance = get_asset_balance(balances, quote_asset)
        available_quote = quote_balance["available"]
        order_size = (available_quote * ORDER_SIZE_PERCENTAGE / 100) / price
    else:  # SELL
        base_asset = pair[:3]  # e.g., XBT in XBTZAR
        base_balance = get_asset_balance(balances, base_asset)
        available_base = base_balance["available"]
        order_size = available_base * ORDER_SIZE_PERCENTAGE / 100

    return order_size

def get_market_making_prices(ticker):
    """Calculate bid and ask prices based on ticker information."""
    bid = float(ticker["bid"])
    ask = float(ticker["ask"])
    mid = (bid + ask) / 2
    spread = mid * SPREAD_PERCENTAGE / 100

    new_bid = mid - spread
    new_ask = mid + spread

    return {
        "bid": round(new_bid, PRICE_PRECISION),
        "ask": round(new_ask, PRICE_PRECISION)
    }

def place_market_making_orders():
    """Place market making orders for all configured trading pairs."""
    balances_response = get_balances()
    if not balances_response:
        log("Failed to get account balances")
        return

    balances = balances_response
    active_orders = {}

    for pair in TRADING_PAIRS:
        log(f"Processing pair: {pair}")

        # Get current market data
        ticker = get_ticker(pair)
        if not ticker:
            log(f"Failed to get ticker for {pair}")
            continue

        # Calculate prices for bid and ask
        prices = get_market_making_prices(ticker)

        # Cancel existing orders for this pair to refresh them
        cancel_all_orders(pair)

        # Place new buy order
        buy_volume = calculate_order_size(pair, balances, prices["bid"], "BUY")
        if buy_volume > 0:
            buy_order_id = create_order(pair, "BID", buy_volume, prices["bid"])
            if buy_order_id:
                active_orders[buy_order_id] = {
                    "pair": pair,
                    "type": "BID",
                    "price": prices["bid"],
                    "volume": buy_volume,
                    "created_at": time.time()
                }

        # Place new sell order
        sell_volume = calculate_order_size(pair, balances, prices["ask"], "ASK")
        if sell_volume > 0:
            sell_order_id = create_order(pair, "ASK", sell_volume, prices["ask"])
            if sell_order_id:
                active_orders[sell_order_id] = {
                    "pair": pair,
                    "type": "ASK",
                    "price": prices["ask"],
                    "volume": sell_volume,
                    "created_at": time.time()
                }

    return active_orders

def main():
    """Main function to run the market maker bot."""
    log("Starting Luno Market Maker Bot")

    if not API_KEY or not API_SECRET:
        log("API_KEY and API_SECRET must be set as environment variables")
        return

    try:
        while True:
            log("Placing market making orders...")
            active_orders = place_market_making_orders()

            # Wait before refreshing orders
            log(f"Waiting {ORDER_EXPIRY_MINUTES} minutes before refreshing orders...")
            time.sleep(ORDER_EXPIRY_MINUTES * 60)
    except KeyboardInterrupt:
        log("Bot stopped by user")
        # Cancel all orders on exit
        cancel_all_orders()
    except Exception as e:
        log(f"Unexpected error: {e}")
        # Cancel all orders on error
        cancel_all_orders()

if __name__ == "__main__":
    main()
