import axios from 'axios';
import { API_ROUTES } from '../../../config';

function prepareData(data, eventType) {
  if (Array.isArray(data)) {
    const requestData = [];
    data.forEach(dataPortion => requestData.concat(prepareData(dataPortion, eventType)));
    return requestData;
  }
  const preparedData = {
    author: data.author,
    id: data.id,
    connector: 'toto',
    connector_name: 'toto',
    event_type: eventType,
    source_type: 'resource',
    component: 'localhost',
    state: Object.prototype.hasOwnProperty.call(data, 'state') ? data.state : 0,
    crecord_type: eventType,
    timestamp: Date.now(),
    resource: data.resource,
    ref_rk: `${data.resource}/${data.component}`,
  };
  if (eventType !== 'snooze') {
    preparedData.state_type = data.state_type ? data.state_type : 1;
  }
  Object.keys(data.customAttributes).forEach((attribute) => {
    preparedData[attribute] = data.customAttributes[attribute];
  });
  return [preparedData];
}

export default {
  namespaced: true,
  actions: {
    async cancelAck(context, data) {
      try {
        await axios.post(API_ROUTES.event, {
          event: prepareData(data, 'ackremove'),
        });
      } catch (e) {
        console.log(e);
      }
    },
    async ack(context, data) {
      try {
        await axios.post(API_ROUTES.event, {
          event: prepareData(data, 'ack'),
        });
      } catch (e) {
        console.log(e);
      }
    },
    async declare(context, data) {
      try {
        await axios.post(API_ROUTES.event, {
          event: prepareData(data, 'declareticket'),
        });
      } catch (e) {
        console.log(e);
      }
    },
    async changeState(context, data) {
      try {
        await axios.post(API_ROUTES.event, {
          event: prepareData(data, 'changestate'),
        });
      } catch (e) {
        console.log(e);
      }
    },
    async cancelAlarm(context, data) {
      try {
        await axios.post(API_ROUTES.event, {
          event: prepareData(data, 'cancel'),
        });
      } catch (e) {
        console.log(e);
      }
    },
    async snooze(context, data) {
      try {
        axios.post(API_ROUTES.event, {
          event: prepareData({
            ...data,
            state: 3,
          }, 'snooze'),
        });
      } catch (e) {
        console.log(e);
      }
    },
  },
};
