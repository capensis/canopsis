import axios from 'axios/index';
import { API_HOST } from '../../../config';

export default {
  namespaced: true,
  actions: {
    cancel: async (context, { comment, alarmData }) => axios.post(`${API_HOST}/event`, {
      event: [
        {
          author: 'root',
          id: 'ac4f92ea-4eda-11e8-841e-0242ac12000a',
          connector: alarmData.connector,
          connector_name: alarmData.connector_name,
          event_type: 'ackremove',
          source_type: 'resource',
          component: alarmData.component,
          state: alarmData.state,
          state_type: alarmData.state_type,
          crecord_type: 'ackremove',
          timestamp: Date.now(),
          resource: alarmData.resource,
          output: comment,
          ref_rk: `${alarmData.resource}/${alarmData.component}`,
        },
      ],
    }),
  },
};
