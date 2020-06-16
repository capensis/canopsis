import axios from 'axios';

/**
 * Active asiox sources
 *
 * @type {{}}
 */
const activeSources = {};

/**
 * Hook for cancelling of the requests by axios source
 *
 * @param {Function} action
 * @param {string} key
 * @return {Promise<void>}
 */
export async function useRequestCancelling(action, key) {
  try {
    const source = axios.CancelToken.source();

    if (activeSources[key]) {
      activeSources[key].cancel();
    }

    activeSources[key] = source;

    await action(source);

    delete activeSources[key];
  } catch (err) {
    if (!axios.isCancel(err)) {
      throw err;
    }
  }
}

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
  if (responseWithError.response && responseWithError.response.status === 401) {
    window.location.reload();
  }
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
