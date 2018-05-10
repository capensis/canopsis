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
      state: data.state ? data.state : 0,
      state_type: data.state_type,
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
        state: dataPortion.state ? dataPortion.state : 0,
        state_type: dataPortion.state_type,
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
    async ack(context, data) {
      return axios.post(`${API_HOST}/event`, {
        event: prepareData(data, 'ack'),
      });
    },
    async declare(context, data) {
      return axios.post(`${API_HOST}/event`, {
        event: prepareData(data, 'declareticket'),
      });
    },
    async snooze(context, data) {
      const requestData = [];
      if (typeof data === 'object') {
        requestData.push({
          author: 'root',
          id: data.id,
          connector: 'toto',
          connector_name: 'toto',
          event_type: 'snooze',
          source_type: 'resource',
          component: 'localhost',
          state: 3,
          crecord_type: 'snooze',
          timestamp: Date.now(),
          resource: data.resource,
          duration: data.customAttributes.duration,
        });
      } else if (Array.isArray(data)) {
        data.forEach((dataPortion) => {
          requestData.push({
            author: 'root',
            id: dataPortion.id,
            connector: 'toto',
            connector_name: 'toto',
            event_type: 'snooze',
            source_type: 'resource',
            component: 'localhost',
            state: 3,
            crecord_type: 'snooze',
            timestamp: Date.now(),
            resource: dataPortion.resource,
            duration: dataPortion.customAttributes.duration,
          });
        });
      }


      return axios.post(`${API_HOST}/event`, {
        event: requestData,
      });
    },
  },
};
