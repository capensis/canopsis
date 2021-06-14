import axios from 'axios';
import { get } from 'lodash';
import Cookies from 'js-cookie';

import { API_BASE_URL, COOKIE_SESSION_KEY } from '@/config';

/**
 * Active axios sources
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
      delete activeSources[key];

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
  if (get(response, 'data.errors', []).length) {
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
  if (responseWithError.response) {
    /**
     * When we will receive 502 or 401 error we must remove cookie to avoid getting a infinity page refreshing
     */
    if ([502, 401].includes(responseWithError.response.status)) {
      Cookies.remove(COOKIE_SESSION_KEY);
      window.location.reload();
    }

    if (responseWithError.response.data) {
      if (responseWithError.response.data.errors) {
        return Promise.reject(responseWithError.response.data.errors);
      }

      return Promise.reject(responseWithError.response.data);
    }
  }

  return Promise.reject(responseWithError);
}

const request = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

request.interceptors.response.use(successResponseHandler, errorResponseHandler);

export default request;
