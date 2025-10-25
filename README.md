# Currency Exchanger

WHY?
> "Because I want. And it's quite funny." — Dalarin1

## Quick Start

### Prerequisites

- Go 1.19 or higher
- An API key from [ExchangeRate-API](https://www.exchangerate-api.com/)

### Installation

1. **Clone & Enter**

   ```bash
   git clone https://github.com/yourusername/currency-exchanger.git
   cd currency-exchanger
   ```

2. **Grab the Dependencies**

    ```bash
    go get github.com/joho/godotenv
    ```

3. **Set Up Your API Key**

    **Option A**: Hardcode it here:

    ```go
    var API_KEY string = "YOUR_API_CODE"

    ```

    **Option B**: Create .env file in project root:

    ```text
    EXCHANGERATE_API_KEY=your_actual_key_here
    ```

    Why download _godotenv_ even for hardcoded version? Because I can.
4. **Build It**

    ```bash
    go build -o currency main.go
    ```

5. **Enjoy this shit**

    ```bash
    ./currency -f USD -t EUR
    ```

## 🎯 Usage

Check USD to EUR:

```bash
./currency -f USD -t EUR
./currency -f USD -e
./currency -f USD -a 100500 -t EUR
```

### Command-line Flags

- `-f` Base currency (default: "RUB")
- `-t` Target currency
- `-a` Amount of base currency to convert to target currency (will break without target currency)
- `-hist` Show historical data about base currency
- `-e` Show enriched data about base currency
- `-h` Show help (there is no help, actually.)
