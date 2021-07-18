# CryptoTax
Tracking cost for taxes is tedious. Let's automate it!

# Disclaimer
The software comes as is without any guarantees or warranties. The author(s) do not claim neither legal nor accounting expertise and will not be liable for the use or misuse of the software. Use at your own risk.

# Canada Requirements
You should maintain the following records on your cryptocurrency transactions:
- the date of the transactions
- the receipts of purchase or transfer of cryptocurrency
- the value of the cryptocurrency in Canadian dollars at the time of the transaction
- the digital wallet records and cryptocurrency addresses
- a description of the transaction and the other party (even if it is just their cryptocurrency address)
- the exchange records
- accounting and legal costs
- the software costs related to managing your tax affairs.

If you are a miner, also keep the following records:
- receipts for the purchase of cryptocurrency mining hardware
- receipts to support your expenses and other records associated with the mining operation (such as power costs, mining pool fees, hardware specifications, maintenance costs, and hardware operation time)
- the mining pool details and records

# Assumptions
- User maintains their records and only uses this software to determine their taxable amount at the end of the year.
- User is in Canada and pays canadian taxes.

# Design
## Data Retrieval
Data can be retrieved in 1 of 2 ways:
1. CSV - named with the service.
2. API token

## Data Structures
Each event should internally be stored as:
- date
- transaction type: buy, sell, exchange, transfer, mining deposit, mining fee
- action (debit)
    - address
    - type
    - amount
    - cadValue
- action (credit)

Note: Fees and exchange rates are encoded
- fees = debit value - credit value
- rate = amount / type

ACB can be calculated by dividing the sum of buy/mining events' credit amounts over the total costs:
E.g.
1. bought 1 BTC for $40000 today with $4 purchase fee:
    - {1, buy, {null, CAD, 40004, 40004}, {btcAddr, BTC, 1, 40000}}
2. bought 1 BTC for $60000 tomorrow with $6 purchase fee:
    - {2, buy, {null, CAD, 60006, 60006}, {btcAddr, BTC, 1, 60000}}
ACB = total cad cost / btc obtained = 100010 / 2 = 50005
3. mined 1 BTC and paid .1 BTC in mining fees at an exchange rate of 50000/btc:
    - {3, mining deposit, null, {btcAddr, BTC, 1, 50000}}
    - {3, mining fee, {btcAddr, BTC, .1, 5000}}
New ACB = SUM(40004, 60006, 5000)/3 = 35003.333333333336