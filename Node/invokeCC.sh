# creating bank:
node invoke bankcc writeBankInfo 1bank kvb chennai 40A 2333s673sxx78 sdr3cfgtdui3 23rfs6vhj148b 897vhessety 86zs0lhtd

# creating business:
node invoke businesscc putNewBusinessInfo 1bus tata 12348901 4000000 23hhnx56s673sxx78 sdr32123d3 23rfs148b 12.4 8.09 0 1000000


# creating wallets
node invoke walletcc newWallet 23hhnx56s673sxx78 1000
node invoke walletcc newWallet 23rfs148b 1000
node invoke walletcc newWallet sdr32123d3 1000
node invoke walletcc newWallet 2333s673sxx78 1000
node invoke walletcc newWallet sdr3cfgtdui3 1000
node invoke walletcc newWallet 23rfs6vhj148b 1000
node invoke walletcc newWallet 897vhessety 1000
node invoke walletcc newWallet 86zs0lhtd 1000

# invoking disbursement transaction
# node invoke txncc newTxnInfo 1txn disbursement 23/04/2018 1loan 1inst 300 1bank 1bus pragadeesh v7b9h
