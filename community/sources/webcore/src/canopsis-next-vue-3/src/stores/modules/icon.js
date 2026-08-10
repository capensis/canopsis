import { keyBy } from 'lodash';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

const convertObjectToFormData = (data = {}) => {
  const formData = new FormData();

  Object.entries(data).forEach(([key, value]) => {
    if (value !== undefined) {
      formData.append(key, value);
    }
  });

  return formData;
};

export const useIconStore = defineStore('icon', () => {
  const registeredIconsById = ref({});

  const fetchListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.icons, { params });

  const fetchItemWithoutStore = ({ id } = {}) =>
    request.get(`${API_ROUTES.icons}/${id}`, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });

  const create = ({ data } = {}) =>
    request.post(API_ROUTES.icons, convertObjectToFormData(data), {
      headers: { 'Content-Type': 'multipart/form-data' },
    });

  const update = ({ id, data } = {}) =>
    request.patch(`${API_ROUTES.icons}/${id}`, convertObjectToFormData(data), {
      headers: { 'Content-Type': 'multipart/form-data' },
    });

  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.icons}/${id}`);

  const addRegisteredIcon = ({ id, icon } = {}) => {
    registeredIconsById.value[id] = icon;
  };

  const removeRegisteredIcon = ({ id } = {}) => {
    delete registeredIconsById.value[id];
  };

  const setRegisteredIcons = ({ icons = [] } = {}) => {
    registeredIconsById.value = keyBy(icons, '_id');
  };

  return {
    registeredIconsById,

    fetchListWithoutStore,
    fetchItemWithoutStore,
    create,
    update,
    remove,
    addRegisteredIcon,
    removeRegisteredIcon,
    setRegisteredIcons,
  };
});
