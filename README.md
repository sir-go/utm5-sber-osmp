## OSMP payments accept server
___
### What it does

This is a server for accepting payments from subscribers via 
the Sberbank-Online service with OSMP protocol.

 - `check` endpoint checks if an account exists in the billing 
   and returns an XML answer with account info
 - `payment` endpoint 
   - checks payment parameters
   - makes  payment in the billing system
   - sends a task for fiscal cheque issue to the pos terminal
   - returns an XML answer
___
### Configure

Flag `-c` sets the configuration file path (default `./config.toml`).

```toml
[service]
    host        = "localhost"               # service address
    port        = 8483                      # service port
    location    = "Europe/Moscow"           # timezone

    [service.timeouts]                      # API handling timeouts
        write   = "10s"
        read    = "10s"
        idle    = "30s"

    [service.users]                         # basic auth simple tokens
        "sber-osmp" = ""
        "test"      = ""

[billing]                                   # billing API credentials
    api_url = ""                            # API URL
    username = ""                           # billing system account username
    password = ""                           # billing system account password
    payment_method = 2                      # payments method code (2 - banks transfert)
    payment_back_method = 8                 # payment return method code (8 - money-back)
    payment_report_retro = "1h"             # the time interval for seeking previous payments

   [billing.tih]                            # billing systems array
       api_prefix = "tih"                   # billing system API prefix
       pay_id_prefix = 352120               # prefix for payments ID
       office = "Тихорецкий филиал"         # office name (for print on cheques
   
   # ...

[osmp]                                      
    check_info = '''Оплата за Интернет'''   # payment title
    time_layout = "02.01.2006_15:04:05"     # payment time layout
    id_max_len = 18                         # maximum lenght for payment ID
    pay_amount_min = 1                      # minimal payment value
    pay_amount_max = 100000                 # maximum payment value

[atol]
    api_url = ""                            # POS terminal API URL

    # POS terminal API request template
    request_template = '''
    {
        "uuid": "%s",
        "request": {
            "type": "sell",
            "items": [{
                "tax": {"type": "none"},
                "name": "абон. плата за Интернет [%d]",
                "type": "position",
                "price": %.2f,
                "amount": %.2f,
                "quantity": 1
            }],
            "total": %.2f,
            "operator": {"name": "%s"},
            "payments": [{"sum": %.2f, "type": "electronically"}],
            "clientInfo": {"emailOrPhone": "%s"},
            "paymentsPlace": "%s",
            "electronically": true
        }
    }
    '''
```
___
### Build & run

```bash
go mod download
go build -o osmp ./cmd/osmp
```
