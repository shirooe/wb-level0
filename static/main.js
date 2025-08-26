const form = document.querySelector('form');
const result = document.querySelector('#result');

form.addEventListener('submit', (e) => {
  e.preventDefault();
  
  const value = document.querySelector('input').value.trim();

  const message = document.querySelector('#message');
  message.classList.remove('success');
  message.classList.remove('error');
  message.classList.add('loading');
  message.innerHTML = 'Loading...';
  result.innerHTML = '';

  fetch(`http://localhost:3000/order/${value}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  })
  .then((response) => {
    if (!response.ok) {
      throw new Error(response.statusText);
    }
    return response.json();
  })
  .then((order) => {
    const html = `
        <div class="result-message"></div>
        <div>
          <h2 class="result-title">Order Summary</h2>
          <p class="result-text">
            <strong>Order UID:</strong> ${order.order_uid}<br>
            <strong>Track Number:</strong> ${order.track_number}<br>
            <strong>Customer ID:</strong> ${order.customer_id}<br>
            <strong>Locale:</strong> ${order.locale}<br>
            <strong>Delivery Service:</strong> ${order.delivery_service}<br>
            <strong>Total Cost:</strong> ${order.payment.amount}<br>
          </p>
          <div class="result-info">
            <div class="payment">
                <h3 class="result-subtitle">Payment Information</h3>
                <p class="result-text">
                  <strong>Transaction:</strong> ${order.payment.transaction}<br>
                  <strong>Request ID:</strong> ${order.payment.request_id}<br>
                  <strong>Currency:</strong> ${order.payment.currency}<br>
                  <strong>Amount:</strong> ${order.payment.amount}<br>
                  <strong>Bank:</strong> ${order.payment.bank}<br>
                  <strong>Delivery Cost:</strong> ${order.payment.delivery_cost}<br>
                  <strong>Goods Total:</strong> ${order.payment.goods_total}<br>
                  <strong>Custom Fee:</strong> ${order.payment.custom_fee}<br>
                </p>
              </div>
            <div class="delivery">
              <h3 class="result-subtitle">Delivery Information</h3>
              <p class="result-text">
                <strong>Name:</strong> ${order.delivery.name}<br>
                <strong>Phone:</strong> ${order.delivery.phone}<br>
                <strong>Zip:</strong> ${order.delivery.zip}<br>
                <strong>City:</strong> ${order.delivery.city}<br>
                <strong>Address:</strong> ${order.delivery.address}<br>
              </p>
            </div>
          </div>
          <h3 class="result-subtitle">Items</h3>
          <ul class="result-list">
            ${order.items.map((item) => `
              <li>
                <strong>Item ID:</strong> ${item.chrt_id}<br>
                <strong>Item Name:</strong> ${item.name}<br>
                <strong>Size:</strong> ${item.size}<br>
                <strong>Brands:</strong> ${item.brand}<br>
                <strong>Price:</strong> ${item.price}<br>
                <strong>Status:</strong> ${item.status}<br>
                <strong>Total Price:</strong> ${item.total_price}<br>
              </li>
            `).join('')}
          </ul>
        </div>
    `;
    result.innerHTML = html;
    message.classList.remove('loading');
    message.classList.add('success');
    message.innerHTML = `Success! Order ID: ${order.order_uid}`;
  })
  .catch((error) => {
    message.classList.remove('loading');
    message.classList.add('error');
    message.innerHTML = `Error: ${error.message}`;
  });
});