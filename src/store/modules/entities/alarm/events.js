import axios from 'axios';
import { API_HOST } from '../../../../config';

function prepareData(id, resource, eventType, additionalAttributes) {
  const commonData = {
    author: 'root',
    id,
    connector: 'toto',
    connector_name: 'toto',
    event_type: eventType,
    source_type: 'resource',
    component: 'localhost',
    state: 0,
    state_type: 1,
    crecord_type: eventType,
    timestamp: Date.now(),
    resource,
    ref_rk: `${resource}/localhost`,
  };

  Object.getOwnPropertyNames(additionalAttributes).forEach((attribute) => {
    commonData[attribute] = additionalAttributes[attribute];
  });

  return commonData;
}

export default {
  namespaced: true,
  actions: {
    async cancel(context, { id, resource, comment }) {
      return axios.post(`${API_HOST}/event`, prepareData(id, resource, 'ackremove', { output: comment }));
    },
  },
};
