import { API_ROUTES } from '@/config';
import { REQUEST_METHODS } from '@/constants';

import request from '@/services/request';

import { createWidgetModule } from '@/store/plugins/entities';

import { prepareAvailabilitiesResponse } from '@/helpers/entities/availability/entity';

export default createWidgetModule({
  route: API_ROUTES.metrics.availability,
  method: REQUEST_METHODS.get,
  dataPreparer: prepareAvailabilitiesResponse,
}, {
  actions: {
    async fetchAvailabilityWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.entityAggregateAvailability, { params });
    },

    fetchAvailabilityHistoryWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.entityAvailability, { params });
    },

    fetchAvailabilityHistoryExport(context, { params }) {
      return request.get(API_ROUTES.metrics.exportAvailabilityByEntity, { params });
    },

    createAvailabilityExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportAvailability, data);
    },

    fetchAvailabilityExport(context, { params, id }) {
      return request.get(`${API_ROUTES.metrics.exportMetric}/${id}`, { params });
    },
  },
});
