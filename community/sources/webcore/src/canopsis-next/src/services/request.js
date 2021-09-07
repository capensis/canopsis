import axios from 'axios';
import { get } from 'lodash';

import { API_HOST, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';

import localStorageService from '@/services/local-storage';

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
 * Prepare axios config before request sending
 *
 * @param {Object} config
 * @returns {*}
 */
function requestHandler(config) {
  if (localStorageService.has(LOCAL_STORAGE_ACCESS_TOKEN_KEY) && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY)}`;
  }

  return config;
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
    const { response, config } = responseWithError;

    /**
     * When we will receive 502 or 401 error we must remove cookie to avoid getting a infinity page refreshing
     */
    if ([502, 401].includes(response.status)) {
      localStorageService.clear();
      window.location.reload();
    }

    if (config.fullResponse) {
      return Promise.reject(response);
    }

    if (response.data) {
      if (response.data.errors) {
        return Promise.reject(response.data.errors);
      }

      return Promise.reject(response.data);
    }
  }

  return Promise.reject(responseWithError);
}

const request = axios.create({ baseURL: API_HOST });

request.interceptors.request.use(requestHandler);
request.interceptors.response.use(successResponseHandler, errorResponseHandler);

export default request;
