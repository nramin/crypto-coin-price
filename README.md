# Crypto Coin Price

Portfolio can be put in crypto.yaml along with your coinmarketcap API key

Under positions field in yaml file is the UCID or coinmarketcap coin id found on each coin's page. The coin id is the key and the value is the quantity in your portfolio for example.

As an example, in the yaml file BTC is included as #1 and ETH is included as #1027.

Run with the following:
```
# go build
# ./crypto-coin-price
```