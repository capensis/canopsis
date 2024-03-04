import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createWidgetModule } from '@/store/plugins/entities';

export default createWidgetModule({ route: API_ROUTES.availability.list }, {
  actions: {
    async fetchAvailabilityWithoutStore(context, { id, params } = {}) {
      // return request.get(API_ROUTES.entityAvailability, { params });
      // eslint-disable-next-line no-console
      console.info('fetchAvailabilityWithoutStore', id, params);
      /**
       * TODO: Should be replaced on real fetch function
       */
      await new Promise(r => setTimeout(r, 2000));

      const minDate = new Date();
      minDate.setDate(minDate.getDate() - 3);

      return {
        availability: {
          uptime: Math.round(Math.random() * 100000),
          downtime: Math.round(Math.random() * 100000),
          inactive_time: Math.round(Math.random() * 1000),
        },
        min_date: Math.round(minDate.getTime() / 1000),
      };
    },

    async fetchAvailabilityHistoryWithoutStore(context, { id, params } = {}) {
      // return request.get(API_ROUTES.entityAvailability, { params });
      // eslint-disable-next-line no-console
      console.info('fetchAvailabilityHistoryWithoutStore', id, params);
      /**
       * TODO: Should be replaced on real fetch function
       */
      await new Promise(r => setTimeout(r, 2000));

      const minDate = new Date();
      minDate.setDate(minDate.getDate() - 3);

      return {
        data: [
          {
            timestamp: new Date(2024, 2, 4, 8).getTime() / 1000,
            uptime: Math.round(Math.random() * 10000),
            downtime: Math.round(Math.random() * 10000),
            inactive_time: Math.round(Math.random() * 1000),
          },
          {
            timestamp: new Date(2024, 2, 4, 9).getTime() / 1000,
            uptime: Math.round(Math.random() * 10000),
            downtime: Math.round(Math.random() * 10000),
            inactive_time: Math.round(Math.random() * 1000),
          },
          {
            timestamp: new Date(2024, 2, 4, 10).getTime() / 1000,
            uptime: Math.round(Math.random() * 10000),
            downtime: Math.round(Math.random() * 10000),
            inactive_time: Math.round(Math.random() * 1000),
          },
          {
            timestamp: new Date(2024, 2, 4, 11).getTime() / 1000,
            uptime: Math.round(Math.random() * 10000),
            downtime: Math.round(Math.random() * 10000),
            inactive_time: Math.round(Math.random() * 1000),
          },
          {
            timestamp: new Date(2024, 2, 4, 12).getTime() / 1000,
            uptime: Math.round(Math.random() * 10000),
            downtime: Math.round(Math.random() * 10000),
            inactive_time: Math.round(Math.random() * 1000),
          },
        ],
      };
    },

    createAvailabilityExport(context, { data }) {
      return request.post(API_ROUTES.availability.exportList, data);
    },

    fetchAvailabilityExport(context, { params, id }) {
      return request.get(`${API_ROUTES.availability.exportList}/${id}`, { params });
    },
  },
});
