import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useEntityCategoryStore = defineStore('entity-category', () => {
  const fetchCategoriesListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.entityCategories, { params });

  return {
    fetchCategoriesListWithoutStore,
  };
});
