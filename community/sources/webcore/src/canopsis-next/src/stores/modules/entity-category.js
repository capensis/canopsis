import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const createEntityCategory = ({ data } = {}) =>
  request.post(API_ROUTES.entityCategories, data);

export const useEntityCategoryStore = defineStore('entity-category', () => {
  const fetchCategoriesListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityCategories, {
      params,
      /**
       * Unlike other list endpoints, entity categories bind repeated `ids` without brackets.
       **/
      paramsSerializer: { indexes: null },
    });

  return {
    create: createEntityCategory,
    fetchCategoriesListWithoutStore,
  };
});
