# OSMP payments accept server
[![Go](https://github.com/sir-go/utm5-sber-osmp/actions/workflows/go.yml/badge.svg)](https://github.com/sir-go/utm5-sber-osmp/actions/workflows/go.yml)

## What it does
This is a service for accepting payments from subscribers via 
the Sberbank-Online service with OSMP protocol.

 - `check` endpoint checks if an account exists in the billing 
   and returns an XML answer with account info
 - `payment` endpoint 
   - checks payment parameters
   - checks if payment already exists in the billing by ID, payment sum and time
   - makes  payment in the billing system
   - sends a task for fiscal cheque issue to the pos terminal
   - returns an XML answer

## Requests parameters
- `action` - what to do: `check` or `payment`
  - `action=check`
    - `account` - subscriber account ID
  - `action=payment`
    - `account` - subscriber account ID
    - `amount` - payment sum
    - `pay_id` - bank internal payment ID
    - `pay_date` - payment time
    - `contact` - (optional) contact e-mail of the payer

## Response XML fields
### check action
- `CODE` - request status (!= 0 if errors are acquired)
- `MESSAGE` - error additional message
- `FIO` - subscriber's full name
- `ADDRESS` - subscriber's home address
- `BALANCE` - subscriber's account balance
- `REC_SUM` - recommended payment sum
- `INFO` - some additional subscriber's information

### payment action
- `CODE` - request status (!= 0 if errors are acquired)
- `MESSAGE` - error additional message
- `EXT_ID` - billing internal payment ID
- `REG_DATE` - payment registration time
- `AMOUNT` - payment registered sum

## Configure
Flag `-c` sets the configuration file path (default `./config.toml`).
Example config - `config.toml`


## Docker
```bash
docker build -t osmp .
docker run --rm -it -v ${PWD}/config.toml:/config.toml:ro osmp:latest
```

## Tests
```bash
go test -v ./...
gosec ./...
```

## Build & run
```bash
go mod download
go build -o osmp ./cmd/main
```
