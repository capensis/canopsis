import axios from 'axios/index';
import { API_HOST } from '../../../config';

export default {
  namespaced: true,
  actions: {
    cancelConfirmation: (context, { comment, alarmData }) => new Promise((resolve, reject) => {
      const testVar;

      axios.post(`${API_HOST}/event`, {
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
      }, {
        headers: {
          Cookie: 'beaker.session.id=3ca4d8435795b8b09204c8f865013ba651e0d2adba0dd7c0f897ac0d265c46192145a5b8',
        },
      })
        .then(() => {
          resolve();
        })
        .catch(reject);
    }),
  },
};
