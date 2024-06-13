import { cloneDeep } from 'lodash';
import Vue from 'vue';
import Faker from 'faker';

import { fakeMeta, fakeParams } from '@unit/data/request-data';

import SetSeveralPlugin from '@/plugins/set-several';

import activeWidgetsModule, { localGetters, localTypes } from '@/store/modules/active-view/active-widgets';

/* eslint-disable-next-line import/no-named-as-default-member */
const { state: initialState, mutations, getters } = activeWidgetsModule;

Vue.use(SetSeveralPlugin);

const mockData = {
  meta: fakeMeta(),
  params: fakeParams(),
  widgetId: Faker.datatype.number(),
  allIds: Faker.datatype.array(5),
};

describe('Active widgets module', () => {
  it('Mutate state after commit FETCH_LIST', () => {
    const { params, widgetId } = mockData;
    const state = cloneDeep(initialState);

    const fetchList = mutations[localTypes.FETCH_LIST];

    fetchList(state, { widgetId, params });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { pending: true, fetchingParams: params },
      },
    });
  });

  it('Mutate state after commit FETCH_LIST_COMPLETED', () => {
    const { meta, params, allIds, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { pending: true, fetchingParams: params },
      },
    };

    const fetchListCompleted = mutations[localTypes.FETCH_LIST_COMPLETED];

    fetchListCompleted(state, { widgetId, allIds, meta });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { fetchingParams: params, allIds, meta, pending: false },
      },
    });
  });

  it('Mutate state after commit FETCH_LIST_FAILED', () => {
    const { meta, params, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { pending: true, meta, fetchingParams: params },
      },
    };

    const fetchListFailed = mutations[localTypes.FETCH_LIST_FAILED];

    fetchListFailed(state, { widgetId });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { pending: false, fetchingParams: params, meta },
      },
    });
  });

  it('Mutate state after commit CLEAR', () => {
    const { meta, params, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { pending: true, meta, fetchingParams: params },
      },
    };

    const clear = mutations[localTypes.CLEAR];

    clear(state);

    expect(state.widgets).toEqual({});
  });

  it('Get allIds data by widget id. Getter: GET_ALL_IDS_BY_WIDGET_ID', () => {
    const { allIds, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { allIds },
      },
    };

    const data = getters[localGetters.GET_ALL_IDS_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual(allIds);
  });

  it('Get meta data by widget id. Getter: GET_META_BY_WIDGET_ID', () => {
    const { meta, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { meta },
      },
    };

    const data = getters[localGetters.GET_META_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual(meta);
  });

  it('Get pending data by widget id. Getter: GET_PENDING_BY_WIDGET_ID', () => {
    const { widgetId } = mockData;
    const pending = true;
    const state = {
      widgets: {
        [widgetId]: { pending },
      },
    };

    const data = getters[localGetters.GET_PENDING_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual(pending);
  });

  it('Get params data by widget id. Getter: GET_FETCHING_PARAMS_BY_WIDGET_ID', () => {
    const { params, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { fetchingParams: params },
      },
    };

    const data = getters[localGetters.GET_FETCHING_PARAMS_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual(params);
  });

  it('Get allIds data by widget id without data. Getter: GET_ALL_IDS_BY_WIDGET_ID', () => {
    const { widgetId } = mockData;
    const state = {
      widgets: {},
    };

    const data = getters[localGetters.GET_ALL_IDS_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual([]);
  });

  it('Get meta data by widget id without data. Getter: GET_META_BY_WIDGET_ID', () => {
    const { widgetId } = mockData;
    const state = {
      widgets: {},
    };

    const data = getters[localGetters.GET_META_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual({});
  });

  it('Get pending data by widget id without data. Getter: GET_PENDING_BY_WIDGET_ID', () => {
    const { widgetId } = mockData;
    const state = {
      widgets: {},
    };

    const data = getters[localGetters.GET_PENDING_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual(false);
  });

  it('Get fetching params data by widget id without data. Getter: GET_FETCHING_PARAMS_BY_WIDGET_ID', () => {
    const { widgetId } = mockData;
    const state = {
      widgets: {},
    };

    const data = getters[localGetters.GET_FETCHING_PARAMS_BY_WIDGET_ID](state)(widgetId);

    expect(data).toEqual({});
  });
});
