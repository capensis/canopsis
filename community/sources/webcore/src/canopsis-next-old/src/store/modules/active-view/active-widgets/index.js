import Vue from 'vue';

export const localTypes = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  CLEAR: 'CLEAR',
};

export const localGetters = {
  GET_ALL_IDS_BY_WIDGET_ID: 'GET_ALL_IDS_BY_WIDGET_ID',
  GET_META_BY_WIDGET_ID: 'GET_META_BY_WIDGET_ID',
  GET_PENDING_BY_WIDGET_ID: 'GET_PENDING_BY_WIDGET_ID',
  GET_FETCHING_PARAMS_BY_WIDGET_ID: 'GET_FETCHING_PARAMS_BY_WIDGET_ID',
};

export const modulePrefix = 'activeView/activeWidgets';

export const types = {
  FETCH_LIST: `${modulePrefix}/${localTypes.FETCH_LIST}`,
  FETCH_LIST_COMPLETED: `${modulePrefix}/${localTypes.FETCH_LIST_COMPLETED}`,
  FETCH_LIST_FAILED: `${modulePrefix}/${localTypes.FETCH_LIST_FAILED}`,

  CLEAR: `${modulePrefix}/${localTypes.CLEAR}`,
};

export const getters = {
  GET_ALL_IDS_BY_WIDGET_ID: `${modulePrefix}/${localGetters.GET_ALL_IDS_BY_WIDGET_ID}`,
  GET_META_BY_WIDGET_ID: `${modulePrefix}/${localGetters.GET_META_BY_WIDGET_ID}`,
  GET_PENDING_BY_WIDGET_ID: `${modulePrefix}/${localGetters.GET_PENDING_BY_WIDGET_ID}`,
  GET_FETCHING_PARAMS_BY_WIDGET_ID: `${modulePrefix}/${localGetters.GET_FETCHING_PARAMS_BY_WIDGET_ID}`,
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    [localGetters.GET_ALL_IDS_BY_WIDGET_ID]: state => widgetId => state.widgets?.[widgetId]?.allIds ?? [],
    [localGetters.GET_META_BY_WIDGET_ID]: state => widgetId => state.widgets[widgetId]?.meta ?? {},
    [localGetters.GET_PENDING_BY_WIDGET_ID]: state => widgetId => state.widgets[widgetId]?.pending ?? false,
    [localGetters.GET_FETCHING_PARAMS_BY_WIDGET_ID]: state => widgetId => state.widgets[widgetId]?.fetchingParams ?? {},
  },
  mutations: {
    [localTypes.FETCH_LIST](state, { widgetId, params }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, fetchingParams: params });
    },
    [localTypes.FETCH_LIST_COMPLETED](state, { widgetId, allIds, meta }) {
      Vue.setSeveral(state.widgets, widgetId, { allIds, meta, pending: false });
    },
    [localTypes.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: false });
    },
    [localTypes.CLEAR](state) {
      state.widgets = {};
    },
  },
  actions: {
    clear({ commit }) {
      commit(localTypes.CLEAR);
    },
  },
};
