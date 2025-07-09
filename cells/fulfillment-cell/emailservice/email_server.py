#!/usr/bin/env python3

import os
import sys
import time
import logging
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from concurrent import futures
from jinja2 import Template

import grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Email templates
ORDER_CONFIRMATION_TEMPLATE = Template("""
<html>
<body>
<h2>Order Confirmation</h2>
<p>Dear Customer,</p>
<p>Thank you for your order! Here are your order details:</p>

<h3>Order ID: {{ order.order_id }}</h3>
<h3>Shipping Information:</h3>
<p>
  <strong>Tracking ID:</strong> {{ order.shipping_tracking_id }}<br>
  <strong>Shipping Cost:</strong> ${{ "%.2f"|format(order.shipping_cost.units + order.shipping_cost.nanos/1000000000) }}<br>
  <strong>Address:</strong><br>
  {{ order.shipping_address.street_address }}<br>
  {{ order.shipping_address.city }}, {{ order.shipping_address.state }} {{ order.shipping_address.zip_code }}<br>
  {{ order.shipping_address.country }}
</p>

<h3>Items Ordered:</h3>
<ul>
{% for item in order.items %}
  <li>{{ item.item.product_id }} - Quantity: {{ item.item.quantity }} - ${{ "%.2f"|format(item.cost.units + item.cost.nanos/1000000000) }}</li>
{% endfor %}
</ul>

<p>Your order will be shipped to the address provided above.</p>
<p>Thank you for your business!</p>

<p>Best regards,<br>
The Online Boutique Team</p>
</body>
</html>
""")

class EmailService:
    def __init__(self):
        self.smtp_server = os.getenv('SMTP_SERVER', 'localhost')
        self.smtp_port = int(os.getenv('SMTP_PORT', '587'))
        self.smtp_username = os.getenv('SMTP_USERNAME', '')
        self.smtp_password = os.getenv('SMTP_PASSWORD', '')
        self.from_email = os.getenv('FROM_EMAIL', 'noreply@onlineboutique.com')
        
    def send_order_confirmation(self, request, context):
        """Send order confirmation email"""
        try:
            email = request.email
            order = request.order
            
            logger.info(f"Sending order confirmation email to {email}")
            
            # Render email template
            html_content = ORDER_CONFIRMATION_TEMPLATE.render(order=order)
            
            # Create email message
            msg = MIMEMultipart('alternative')
            msg['Subject'] = f"Order Confirmation - {order.order_id}"
            msg['From'] = self.from_email
            msg['To'] = email
            
            # Add HTML content
            html_part = MIMEText(html_content, 'html')
            msg.attach(html_part)
            
            # Send email (mock implementation)
            self._send_email(msg)
            
            logger.info(f"Order confirmation email sent successfully to {email}")
            return {}  # Empty response
            
        except Exception as e:
            logger.error(f"Failed to send order confirmation email: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Failed to send email: {str(e)}")
            return {}
    
    def _send_email(self, msg):
        """Mock email sending - in production this would use actual SMTP"""
        logger.info(f"Mock email sent: {msg['Subject']} to {msg['To']}")
        
        # In a real implementation, you would use:
        # with smtplib.SMTP(self.smtp_server, self.smtp_port) as server:
        #     if self.smtp_username:
        #         server.starttls()
        #         server.login(self.smtp_username, self.smtp_password)
        #     server.send_message(msg)
        
        # For demo purposes, just log the email content
        logger.info(f"Email content: {msg.get_payload()}")

def serve():
    port = os.environ.get('PORT', '8080')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    # Add services here when proto bindings are available
    # email_pb2_grpc.add_EmailServiceServicer_to_server(EmailService(), server)
    # health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)
    
    listen_addr = f'[::]:{port}'
    server.add_insecure_port(listen_addr)
    
    logger.info(f"Starting email service on {listen_addr}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Received interrupt signal, shutting down...")
        server.stop(0)

if __name__ == '__main__':
    serve()