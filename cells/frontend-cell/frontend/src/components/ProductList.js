import React from 'react';

const ProductList = () => {
  const products = [
    { id: '1', name: 'Vintage Typewriter', price: 67.99, image: '/api/placeholder/300/200' },
    { id: '2', name: 'Film Camera', price: 2295.00, image: '/api/placeholder/300/200' },
    { id: '3', name: 'Vintage Record Player', price: 65.50, image: '/api/placeholder/300/200' },
  ];

  return (
    <div className="product-list">
      <h2>Featured Products</h2>
      <div className="products-grid">
        {products.map(product => (
          <div key={product.id} className="product-card">
            <img src={product.image} alt={product.name} />
            <h3>{product.name}</h3>
            <p className="price">${product.price}</p>
            <button className="add-to-cart">Add to Cart</button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ProductList;