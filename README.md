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
Legal requirements:
- the value of the cryptocurrency in Canadian dollars at the time of the transaction
- the digital wallet records and cryptocurrency addresses
- a description of the transaction and the other party (even if it is just their cryptocurrency address)

Each event should internally be stored as:
- timestamp
- coinType
- coinQuantity
- cadValue
- walletAddress
- transactionType (buy, sell, deposit, withdrawal, mining, interest, fee)

## Implicit fees
- transfer fees = total withdrawals - total deposits
