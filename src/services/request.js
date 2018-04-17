import axios from 'axios';

import { API_HOST, NON_AUTH_API_ROUTES, AUTH_TOKEN_KEY } from '@/config';

import localStorageDataSource from './local-storage-data-source';

/**
 * Set auth token in the authorization header under the 'Authorization' key
 *
 * @param {Object} config
 * @return {Object}
 */
function setAuthTokenHeader(config) {
  const isAuthHeaderRequired = NON_AUTH_API_ROUTES.every(url => config.url.indexOf(url) === -1);

  if (isAuthHeaderRequired) {
    config.headers.common.Authorization = `Bearer ${localStorageDataSource.getItem(AUTH_TOKEN_KEY)}`;
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
  if (response.data.errors && response.data.errors.length) {
    return Promise.reject(response.data.errors);
  } else if (response.data.data) {
    return response.data.data;
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

const request = axios.create({ baseURL: API_HOST });

request.interceptors.request.use(setAuthTokenHeader);
request.interceptors.response.use(successResponseHandler, errorResponseHandler);

export default request;
