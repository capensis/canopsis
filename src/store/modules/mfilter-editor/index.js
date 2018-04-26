import parseFilterToRequest from '@/services/mfilter-editor/parseFilterToRequest';
import parseGroupToFilter from '@/services/mfilter-editor/parseRequestToFilter';

const types = {
  CHANGE_TAB: 'CHANGE_TAB',
  UPDATE_FILTER: 'UPDATE_FILTER',
};

export default {
  namespaced: true,

  state: {
    filter: [{
      condition: '$or',
      groups: [],
      rules: [],
    }],
    possibleFields: ['component_name', 'connector_name'],
    activeTab: 0,
  },

  getters: {
    filter: state => state.filter,
    possibleFields: state => state.possibleFields,
    activeTab: state => state.activeTab,
    filter2request: state => parseFilterToRequest(state.filter),
  },

  mutations: {
    [types.CHANGE_TAB](state, payload) {
      state.activeTab = payload;
    },
    [types.UPDATE_FILTER](state, payload) {
      state.filter = [parseGroupToFilter(payload)];
    },
  },

  actions: {
    changeActiveTab(context, activeTab) {
      context.commit(types.CHANGE_TAB, activeTab);
    },
    updateFilter(context, newRequest) {
      context.commit(types.UPDATE_FILTER, newRequest);
    },
  },
};
