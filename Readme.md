# Exchange Rates CLI client


This CLI application is designed to give information about exchange rates.


### Using

1. Clone or download repo. Make sure that you unzipped it if you choosed "download zip".

```bash
cd ./exchange-rates-cli

```

2. Build project

```bash 
go build -o exchange-rates
```

3. Run
```bash
./exchange-rate get --currency EUR 
```

4. Convetion between two foreign currencies
```bash
./exchange-rate convert -s USD -t EUR -a 1000
```