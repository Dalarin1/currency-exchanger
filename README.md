#  Currency Exchanger

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
4. **Build It**
    ```bash
    go build -o currency main.go
    ```
5. **Enjoy this shit**
    ```bash
    ./currency -f USD -t EUR
    ```
    