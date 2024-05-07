export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore() {
      return {
        data: [
          { created: 123123123, count: 1124, status: 0 },
          { created: 123123223, count: 10000, status: 1 },
          { created: 123123423, count: 10000, status: 1 },
        ],
        meta: {
          page: 1,
          total_count: 3,
          status: 0,
        },
      };
    },

    fetchEventsListWithoutStore() {
      return {
        data: [
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
          { timestamp: 123123123, event_type: 'check', source_type: 'resource', connector: 'centreon', connector_name: 'centreon_0_SUP_43', resource: 'check_ping' },
        ],
        meta: {
          page: 1,
          total_count: 10,
          status: 0,
        },
      };
    },

    launch() {
      // return request.post(`${API_ROUTES.eventsRecording}/launch`);
    },

    stop() {
      // return request.post(`${API_ROUTES.eventsRecording}/stop`);
    },

    remove() {
      // return request.remove(`${API_ROUTES.eventsRecording}/${id}`);
    },
  },
};
