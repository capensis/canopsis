import { API_ROUTES } from '@/config';
import { REQUEST_METHODS } from '@/constants';

import request from '@/services/request';

import { createWidgetModule } from '@/store/plugins/entities';

const prepareAvailabilityData = ({ data }) => data.map(availability => ({
  ...availability,
  entity: {
    ...availability.entity,
    category: {
      name: availability.entity.category,
    },
    infos: availability.entity.infos
      ? Object.entries(availability.entity.infos).reduce((acc, [name, value]) => {
        acc[name] = {
          value,
          name,
        };

        return acc;
      }, {})
      : {},
  },
}));

export default createWidgetModule({
  route: API_ROUTES.metrics.availability,
  method: REQUEST_METHODS.get,
  dataPreparer: prepareAvailabilityData,
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
