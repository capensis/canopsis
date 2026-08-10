import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

const convertObjectToFormData = (data = {}) =>
  Object.entries(data).reduce((formData, [key, value]) => {
    formData.append(key, value);

    return formData;
  }, new FormData());

export const usePbehaviorExceptionStore = createCRUDStore({
  id: 'pbehaviorException',
  route: API_ROUTES.pbehavior.exceptions,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    import({ data } = {}) {
      return request.post(API_ROUTES.pbehavior.exceptionImport, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    bulkHide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkExceptionsHide, data);
    },

    bulkUnhide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkExceptionsUnhide, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.pbehavior.bulkExceptions, { data });
    },
  },
});
