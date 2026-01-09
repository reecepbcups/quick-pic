## Summary

```bash
# SQLite backend
just run-sqlite

# Blockchain backend (in separate terminals)
just anvil              # Terminal 1: Start Anvil
just deploy-contract    # Terminal 2: Deploy contract (note the address)
just run-blockchain 0x5FbDB2315678afecb367f032d93F642f64180aa3  # Terminal 2: Run server

# Other helpers
just build-contracts    # Build Solidity contracts
```
