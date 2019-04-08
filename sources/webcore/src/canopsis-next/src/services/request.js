import axios from 'axios';

/**
 * Check error field inside successful response and reject them
 *
 * @param {Object} response
 * @returns {Object}
 */
function successResponseHandler(response) {
  if (response.data.errors && response.data.errors.length) {
    return Promise.reject(response.data.errors);
  }
  return response.data;
}

/**
 * Check error field inside unsuccessful response, map and reject them
 *
 * @param {Object} responseWithError
 * @returns {Object}
 */
function errorResponseHandler(responseWithError) {
  if (responseWithError.response && responseWithError.response.data) {
    if (responseWithError.response.data.errors) {
      return Promise.reject(responseWithError.response.data.errors);
    }

    return Promise.reject(responseWithError.response.data);
  }

  return Promise.reject(responseWithError);
}

const request = axios.create({
  baseURL: process.env.NODE_ENV === 'production' ? '' : '/api',
  withCredentials: true,
});

request.interceptors.response.use(successResponseHandler, errorResponseHandler);

export default request;
