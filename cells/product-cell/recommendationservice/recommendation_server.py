#!/usr/bin/env python3

import os
import random
import time
import logging
from concurrent import futures

import grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Product catalog for recommendations
product_catalog = [
    'OLJCESPC7Z',  # Sunglasses
    '66VCHSJNUP',  # Tank Top
    '1YMWWN1N4O',  # Watch
    'L9ECAV7KIM',  # Loafers
    '2ZYFJ3GM2N',  # Hairdryer
    '0PUK6V6EV0',  # Candle Holder
    'LS4PSXUNUM',  # Salt & Pepper Shakers
    '9SIQT8TOJO',  # Bamboo Glass Jar
    '6E92ZMYYFZ',  # Mug
]

# Simple recommendation logic based on product categories
recommendation_map = {
    'OLJCESPC7Z': ['1YMWWN1N4O', '66VCHSJNUP'],  # Sunglasses -> Watch, Tank Top
    '66VCHSJNUP': ['OLJCESPC7Z', 'L9ECAV7KIM'],  # Tank Top -> Sunglasses, Loafers
    '1YMWWN1N4O': ['OLJCESPC7Z', 'L9ECAV7KIM'],  # Watch -> Sunglasses, Loafers
    'L9ECAV7KIM': ['1YMWWN1N4O', '66VCHSJNUP'],  # Loafers -> Watch, Tank Top
    '2ZYFJ3GM2N': ['0PUK6V6EV0', '6E92ZMYYFZ'],  # Hairdryer -> Candle Holder, Mug
    '0PUK6V6EV0': ['LS4PSXUNUM', '9SIQT8TOJO'],  # Candle Holder -> Salt & Pepper, Bamboo Jar
    'LS4PSXUNUM': ['9SIQT8TOJO', '6E92ZMYYFZ'],  # Salt & Pepper -> Bamboo Jar, Mug
    '9SIQT8TOJO': ['LS4PSXUNUM', '6E92ZMYYFZ'],  # Bamboo Jar -> Salt & Pepper, Mug
    '6E92ZMYYFZ': ['9SIQT8TOJO', 'LS4PSXUNUM'],  # Mug -> Bamboo Jar, Salt & Pepper
}

class RecommendationService:
    def ListRecommendations(self, request, context):
        """Get product recommendations based on user's cart"""
        max_responses = 5
        
        # Filter products that are already in the cart
        product_ids = list(request.product_ids)
        filtered_products = []
        
        # Get recommendations for each product in the cart
        for product_id in product_ids:
            if product_id in recommendation_map:
                filtered_products.extend(recommendation_map[product_id])
        
        # If no specific recommendations, return random products
        if not filtered_products:
            filtered_products = product_catalog.copy()
        
        # Remove duplicates and products already in cart
        filtered_products = list(set(filtered_products))
        filtered_products = [p for p in filtered_products if p not in product_ids]
        
        # Randomize and limit results
        random.shuffle(filtered_products)
        num_products = min(max_responses, len(filtered_products))
        
        # Return response structure (simplified without proto)
        response = {
            'product_ids': filtered_products[:num_products]
        }
        
        logger.info(f"Returning {num_products} recommendations for user {request.user_id}")
        return response

def serve():
    port = os.environ.get('PORT', '8080')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    # Add services here when proto bindings are available
    # recommendation_pb2_grpc.add_RecommendationServiceServicer_to_server(RecommendationService(), server)
    # health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)
    
    listen_addr = f'[::]:{port}'
    server.add_insecure_port(listen_addr)
    
    logger.info(f"Starting recommendation service on {listen_addr}")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()