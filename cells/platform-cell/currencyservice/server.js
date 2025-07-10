const path = require('path');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const pino = require('pino');
const { v4: uuidv4 } = require('uuid');

const PROTO_PATH = path.join(__dirname, 'demo.proto');
const PORT = process.env.PORT || 8081;

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const shopProto = grpc.loadPackageDefinition(packageDefinition).hipstershop;
const logger = pino({ level: 'info' });

// Supported currencies
const supportedCurrencies = [
  'USD', 'EUR', 'GBP', 'JPY', 'CAD', 'AUD', 'CHF', 'CNY', 'SEK', 'NOK', 'DKK', 'PLN', 'CZK', 'HUF', 'RON', 'BGN', 'HRK', 'RUB', 'UAH', 'TRY', 'ILS', 'KRW', 'SGD', 'HKD', 'NZD', 'MXN', 'BRL', 'ARS', 'CLP', 'COP', 'PEN', 'UYU', 'VES', 'INR', 'PKR', 'LKR', 'NPR', 'BDT', 'MYR', 'THB', 'IDR', 'VND', 'PHP', 'KHR', 'LAK', 'MMK', 'BND', 'FJD', 'PGK', 'SBD', 'VUV', 'WST', 'TOP', 'BWP', 'ZAR', 'NAD', 'SZL', 'LSL', 'ZMW', 'MWK', 'TZS', 'UGX', 'KES', 'RWF', 'BIF', 'DJF', 'ETB', 'SOS', 'XAF', 'XOF', 'XPF', 'KMF', 'GNF', 'SLL', 'LRD', 'GMD', 'CVE', 'STD', 'AOA', 'MZN', 'MGF', 'KMF', 'SCR', 'MUR', 'EGP', 'LYD', 'TND', 'DZD', 'MAD', 'GHS', 'NGN', 'XAG', 'XAU', 'XPD', 'XPT',
];

// Mock exchange rates (in production, this would come from a real API)
const exchangeRates = {
  USD: 1.0,
  EUR: 0.85,
  GBP: 0.73,
  JPY: 110.0,
  CAD: 1.25,
  AUD: 1.35,
  CHF: 0.92,
  CNY: 6.45,
  SEK: 8.6,
  NOK: 8.5,
  DKK: 6.36,
  PLN: 3.9,
  CZK: 22.0,
  HUF: 295.0,
  RON: 4.2,
  BGN: 1.66,
  HRK: 6.4,
  RUB: 74.0,
  UAH: 27.0,
  TRY: 8.5,
  ILS: 3.3,
  KRW: 1180.0,
  SGD: 1.35,
  HKD: 7.75,
  NZD: 1.4,
  MXN: 20.0,
  BRL: 5.2,
  ARS: 95.0,
  CLP: 800.0,
  COP: 3700.0,
  PEN: 4.0,
  UYU: 43.0,
  VES: 4.0,
  INR: 74.0,
  PKR: 160.0,
  LKR: 200.0,
  NPR: 118.0,
  BDT: 85.0,
  MYR: 4.15,
  THB: 31.0,
  IDR: 14300.0,
  VND: 23000.0,
  PHP: 50.0,
  KHR: 4080.0,
  LAK: 9500.0,
  MMK: 1400.0,
  BND: 1.35,
  FJD: 2.1,
  PGK: 3.5,
  SBD: 8.0,
  VUV: 112.0,
  WST: 2.6,
  TOP: 2.3,
  BWP: 11.0,
  ZAR: 14.5,
  NAD: 14.5,
  SZL: 14.5,
  LSL: 14.5,
  ZMW: 18.0,
  MWK: 810.0,
  TZS: 2300.0,
  UGX: 3550.0,
  KES: 108.0,
  RWF: 1000.0,
  BIF: 1970.0,
  DJF: 177.0,
  ETB: 44.0,
  SOS: 580.0,
  XAF: 558.0,
  XOF: 558.0,
  XPF: 102.0,
  KMF: 419.0,
  GNF: 10500.0,
  SLL: 10200.0,
  LRD: 170.0,
  GMD: 52.0,
  CVE: 94.0,
  STD: 21000.0,
  AOA: 650.0,
  MZN: 64.0,
  MGF: 4000.0,
  SCR: 13.5,
  MUR: 42.0,
  EGP: 15.7,
  LYD: 4.5,
  TND: 2.8,
  DZD: 135.0,
  MAD: 9.0,
  GHS: 6.0,
  NGN: 410.0,
  XAG: 0.04,
  XAU: 0.0005,
  XPD: 0.0004,
  XPT: 0.001,
};

function getSupportedCurrencies(call, callback) {
  logger.info('Getting supported currencies');
  callback(null, { currency_codes: supportedCurrencies });
}

function convert(call, callback) {
  try {
    const request = call.request;
    logger.info('Converting currency', { from: request.from, to: request.to_code });
    
    const fromCurrency = request.from.currency_code;
    const toCurrency = request.to_code;
    const amount = parseFloat(request.from.units) + parseFloat(request.from.nanos) / 1e9;
    
    // Validate currencies
    if (!supportedCurrencies.includes(fromCurrency)) {
      return callback(new Error(`Unsupported currency: ${fromCurrency}`));
    }
    
    if (!supportedCurrencies.includes(toCurrency)) {
      return callback(new Error(`Unsupported currency: ${toCurrency}`));
    }
    
    // Convert to USD first, then to target currency
    const usdAmount = amount / exchangeRates[fromCurrency];
    const convertedAmount = usdAmount * exchangeRates[toCurrency];
    
    const units = Math.floor(convertedAmount);
    const nanos = Math.round((convertedAmount - units) * 1e9);
    
    const result = {
      currency_code: toCurrency,
      units: units,
      nanos: nanos,
    };
    
    logger.info('Conversion result', { result });
    callback(null, result);
  } catch (error) {
    logger.error('Error during currency conversion', { error });
    callback(error);
  }
}

function main() {
  const server = new grpc.Server();
  
  server.addService(shopProto.CurrencyService.service, {
    getSupportedCurrencies: getSupportedCurrencies,
    convert: convert,
  });
  
  server.bindAsync(
    `0.0.0.0:${PORT}`,
    grpc.ServerCredentials.createInsecure(),
    () => {
      logger.info(`Currency service listening on port ${PORT}`);
      server.start();
    }
  );
}

if (require.main === module) {
  main();
}

module.exports = { main };
