import parseFilterToRequest from '@/services/mfilter-editor/parseFilterToRequest';
import parseGroupToFilter from '@/services/mfilter-editor/parseRequestToFilter';

const types = {
  CHANGE_TAB: 'CHANGE_TAB',
  UPDATE_FILTER: 'UPDATE_FILTER',
  ADD_PARSE_ERROR: 'ADD_PARSE_ERROR',
  DELETE_PARSE_ERROR: 'DELETE_PARSE_ERROR',
};

export default {
  namespaced: true,

  state: {
    filter: [{
      condition: '$or',
      groups: [],
      rules: [],
    }],
    possibleFields: ['component_name', 'connector_name', 'connector', 'resource'],
    activeTab: 0,
    parseError: '',
  },

  getters: {
    filter: state => state.filter,
    possibleFields: state => state.possibleFields,
    activeTab: state => state.activeTab,
    request: (state) => {
      try {
        return parseFilterToRequest(state.filter);
      } catch (e) {
        return e;
      }
    },
    parseError: state => state.parseError,
  },

  mutations: {
    [types.CHANGE_TAB](state, payload) {
      state.activeTab = payload;
    },
    [types.UPDATE_FILTER](state, payload) {
      try {
        const newFilter = parseGroupToFilter(payload);
        state.filter = [newFilter];
      } catch (error) {
        state.parseError = error.message;
      }
    },
    [types.ADD_PARSE_ERROR](state, payload) {
      state.parseError = payload;
    },
    [types.DELETE_PARSE_ERROR](state) {
      state.parseError = '';
    },
  },

  actions: {
    changeActiveTab(context, activeTab) {
      context.commit(types.CHANGE_TAB, activeTab);
    },
    updateFilter(context, newRequest) {
      context.commit(types.UPDATE_FILTER, newRequest);
    },
    onParseError(context, error) {
      context.commit(types.ADD_PARSE_ERROR, error);
    },
    deleteParseError(context) {
      context.commit(types.DELETE_PARSE_ERROR);
    },
  },
};
