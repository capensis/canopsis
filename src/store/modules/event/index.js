import axios from 'axios';
import { API_HOST } from '../../../config';

function prepareData(data, eventType) {
  const requestData = [];

  if (typeof data === 'object') {
    const commonData = {
      author: 'root',
      id: data.id,
      connector: 'toto',
      connector_name: 'toto',
      event_type: eventType,
      source_type: 'resource',
      component: 'localhost',
      state: 0,
      state_type: 1,
      crecord_type: eventType,
      timestamp: Date.now(),
      resource: data.resource,
      ref_rk: `${data.resource}/localhost`,
    };
    Object.getOwnPropertyNames(data.customAttributes).forEach((attribute) => {
      commonData[attribute] = data.customAttributes[attribute];
    });
    requestData.push(commonData);
  } else if (Array.isArray(data)) {
    data.forEach((dataPortion) => {
      const commonData = {
        author: 'root',
        id: dataPortion.id,
        connector: 'toto',
        connector_name: 'toto',
        event_type: eventType,
        source_type: 'resource',
        component: 'localhost',
        state: 0,
        state_type: 1,
        crecord_type: eventType,
        timestamp: Date.now(),
        resource: dataPortion.resource,
        ref_rk: `${dataPortion.resource}/localhost`,
      };
      Object.getOwnPropertyNames(dataPortion.customAttributes).forEach((attribute) => {
        commonData[attribute] = dataPortion.customAttributes[attribute];
      });
      requestData.push(commonData);
    });
  }

  return requestData;
}

export default {
  namespaced: true,
  actions: {
    async cancel(context, data) {
      return axios.post(`${API_HOST}/event`, {
        event: prepareData(data, 'ackremove'),
      });
    },
    async fastAck(context, data) {
      return axios.post(`${API_HOST}/event`, {
        event: prepareData(data, 'ack'),
      });
    },
  },
};
