import { parseGroup2Filter, parseRule2Request, parseGroup2Request } from '@/services/mfilter-editor/parse';

export default {
  namespaced: true,
  state: {
    request: {},
    filter: [{
      condition: '$or',
      groups: [],
      rules: [],
    }],
    possibleFields: ['component_name', 'connector_name'],
    activeTab: 0,
  },
  getters: {
    filter2request(state) {
      const request = {};

      request[state.filter[0].condition] = [];

      state.filter[0].rules.map((rule) => {
        request[state.filter[0].condition].push(parseRule2Request(rule));
        return request;
      });
      state.filter[0].groups.map((group) => {
        request[state.filter[0].condition].push(parseGroup2Request(group));
        return request;
      });

      return request;
    },
  },
  mutations: {
    changeActiveTab(state, payload) {
      state.activeTab = payload;
    },
    updateFilter(state, payload) {
      state.filter = [parseGroup2Filter(payload)];
    },
  },
  actions: {
    changeActiveTab(context, activeTab) {
      context.commit('changeActiveTab', activeTab);
    },
    updateFilter(context, newRequest) {
      context.commit('updateFilter', newRequest);
    },
  },
};
