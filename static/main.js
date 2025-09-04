const form = document.querySelector('#order_form');
const testForm = document.querySelector('#test_form');
const result = document.querySelector('#result');
const message = document.querySelector('#message');

// поиск по order_uid
form.addEventListener('submit', (e) => {
  e.preventDefault();
  
  const value = document.querySelector('input').value.trim();

  messageLoadingState()
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
    messageSuccessState(order.order_uid);
  })
  .catch((error) => {
    messageErrorState(error.message)
  });
});

// отправка тестовой модели заказа (131 строка)
testForm.addEventListener('submit', (e) => {
  e.preventDefault();

  messageLoadingState()

  fetch('http://localhost:3000/order', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(model),
  })
  .then((response) => {
    if (!response.ok) {
      throw new Error(response.statusText);
    }
    return response.ok;
  })
  .then(() => {
    messageSuccessState(model.order_uid);
  })
  .catch((error) => {
    messageErrorState(error.message)
  });
})

// состояние загрузки
function messageLoadingState() {
  message.classList.remove('success');
  message.classList.remove('error');
  message.classList.add('loading');
  message.innerHTML = 'Loading...'; 
}

// состояние успешно
function messageSuccessState(id) {
  message.classList.remove('loading');
  message.classList.add('success');
  message.innerHTML = `Success! Order ID: ${id}`;
}

// состояние ошибки
function messageErrorState(error) {
  message.classList.remove('loading');
  message.classList.add('error');
  message.innerHTML = `Error: ${error}`;
}

const model = {
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}