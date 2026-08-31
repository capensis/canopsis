import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useAssociativeTableStore = defineStore('associativeTable', () => {
  /**
   * Fetch associative table content by name.
   *
   * @param {Object} [params={}]
   * @param {string} [params.name] - Associative table name.
   * @returns {Promise<Object>}
   */
  const fetch = async ({ name } = {}) => {
    const { content } = await request.get(API_ROUTES.associativeTable, {
      params: { name },
    });

    return content;
  };

  /**
   * Create associative table content by name.
   *
   * @param {Object} [params={}]
   * @param {string} [params.name] - Associative table name.
   * @param {Object} [params.data] - Associative table content.
   * @returns {Promise<Object>}
   */
  const create = async ({ name, data } = {}) => {
    const { content } = await request.post(API_ROUTES.associativeTable, {
      name,
      content: data,
    });

    return content;
  };

  /**
   * Update associative table content by name.
   *
   * @param {Object} [payload={}] - Update payload.
   * @returns {Promise<Object>}
   */
  const update = payload => create(payload);

  /**
   * Remove associative table by name.
   *
   * @param {Object} [params={}]
   * @param {string} [params.name] - Associative table name.
   * @returns {Promise<Object>}
   */
  const remove = ({ name } = {}) => request.delete(API_ROUTES.associativeTable, {
      params: { name },
    });

  return {
    fetch,
    create,
    update,
    remove,
  };
});
