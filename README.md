# My exchange script but written in Go



# Simple script that converts between GBP and EUR

## Usage:

```xcgo currency value```

currency: e or p

e: converts Euro to Pound
p: converts Pound to Euro

value: float value to be converted

Example to convert 1000£ to €: ```xcgo p 1000```

The script will check for an existing exchange.json file with up to date information. If the file does not exist or the data is not from today, it will fetch new data. This will prevent calls to the API if the retrieved information is from today. Because the exchange value is calculated at 23:59:59 we need to subtract one day from today.

## Configuration:

These two variables need to be configured:

### api_key_file

This script uses http://currencyapi.com for data. You need to setup the variable api_key_file to point to a file containing the API key you setup on that site.

### json_file

Should point to a json file the script will create if it doesn't exist or read from if it does. This file keeps the last fetched data, to avoid constant calls to the API when the fetched data is still recent (last day)
