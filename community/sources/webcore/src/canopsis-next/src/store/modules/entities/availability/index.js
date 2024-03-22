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
        data: Array.from({ length: 5 }, (_, index) => {
          const uptime = Math.round(Math.random() * 10000);
          const downtime = Math.round(Math.random() * 10000);
          const totalTime = uptime + downtime;
          const uptimeShare = ((uptime / totalTime) * 100).toFixed(2);
          const downtimeShare = 100 - uptimeShare;

          return ({
            timestamp: new Date(2024, 2, 4, 8 + index).getTime() / 1000,

            uptime_duration: uptime,
            downtime_duration: downtime,

            uptime_share: uptimeShare,
            downtime_share: downtimeShare,

            uptime_share_history: uptimeShare + Math.round(Math.random() * 20) - 10,
            downtime_share_history: downtimeShare + Math.round(Math.random() * 20) - 10,
          });
        }),
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
