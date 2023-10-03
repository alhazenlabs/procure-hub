// config.js
const IS_DEV = true
let BASE_URL;

if (IS_DEV) {
  BASE_URL = 'http://localhost:8080'; // Replace 'your-backend-port' with the actual port
} else {
  BASE_URL = ''; // Set the production or default URL here
}

export default BASE_URL;