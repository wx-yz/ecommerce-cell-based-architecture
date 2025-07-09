const express = require('express');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');
const { v4: uuidv4 } = require('uuid');
const pino = require('pino');

const PORT = process.env.PORT || 50051;
const PROTO_PATH = path.join(__dirname, '../../../shared/protos/demo.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const shopProto = grpc.loadPackageDefinition(packageDefinition).hipstershop;
const logger = pino({ level: 'info' });

class PaymentService {
  charge(call, callback) {
    try {
      const request = call.request;
      const { amount, credit_card } = request;
      
      logger.info('Processing payment', { 
        amount: amount, 
        cardNumber: credit_card.credit_card_number.slice(-4) 
      });
      
      // Simulate payment processing
      this.validateCreditCard(credit_card);
      
      // Generate transaction ID
      const transactionId = uuidv4();
      
      // Simulate processing delay
      setTimeout(() => {
        logger.info('Payment processed successfully', { transactionId });
        
        callback(null, {
          transaction_id: transactionId
        });
      }, 100);
      
    } catch (error) {
      logger.error('Payment processing failed', { error: error.message });
      callback(error);
    }
  }
  
  validateCreditCard(creditCard) {
    // Basic validation
    if (!creditCard.credit_card_number || creditCard.credit_card_number.length < 13) {
      throw new Error('Invalid credit card number');
    }
    
    if (!creditCard.credit_card_cvv || creditCard.credit_card_cvv < 100) {
      throw new Error('Invalid CVV');
    }
    
    if (!creditCard.credit_card_expiration_month || 
        creditCard.credit_card_expiration_month < 1 || 
        creditCard.credit_card_expiration_month > 12) {
      throw new Error('Invalid expiration month');
    }
    
    if (!creditCard.credit_card_expiration_year || 
        creditCard.credit_card_expiration_year < new Date().getFullYear()) {
      throw new Error('Credit card expired');
    }
    
    // Simulate declined cards
    if (creditCard.credit_card_number.startsWith('4000000000000002')) {
      throw new Error('Credit card declined');
    }
    
    return true;
  }
}

function main() {
  const server = new grpc.Server();
  
  server.addService(shopProto.PaymentService.service, new PaymentService());
  
  server.bindAsync(
    `0.0.0.0:${PORT}`,
    grpc.ServerCredentials.createInsecure(),
    (error, port) => {
      if (error) {
        logger.error('Failed to bind server', { error });
        return;
      }
      
      logger.info(`Payment service listening on port ${port}`);
      server.start();
    }
  );
  
  // Graceful shutdown
  process.on('SIGTERM', () => {
    logger.info('Received SIGTERM, shutting down gracefully');
    server.tryShutdown(() => {
      logger.info('Server shut down');
      process.exit(0);
    });
  });
}

if (require.main === module) {
  main();
}

module.exports = { PaymentService };