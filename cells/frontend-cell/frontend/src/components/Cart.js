import React from 'react';
import { Link } from 'react-router-dom';

const Cart = () => {
  return (
    <div className="cart">
      <h2>Your Cart</h2>
      <p>Your cart is empty</p>
      <Link to="/">Continue Shopping</Link>
    </div>
  );
};

export default Cart;