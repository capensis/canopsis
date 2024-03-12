import { cloneDeep } from 'lodash';
import Vue from 'vue';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { fakeAlarm, fakeAlarms, fakeAlarmsResponse } from '@unit/data/alarm';
import { fakeMeta, fakeParams } from '@unit/data/request-data';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import SetSeveralPlugin from '@/plugins/set-several';

import request from '@/services/request';

import alarmModule, { types } from '@/store/modules/entities/alarm';

const { actions, state: initialState, mutations, getters } = alarmModule;

Vue.use(SetSeveralPlugin);

const mockData = {
  alarmsResponse: fakeAlarmsResponse({ count: 10 }),
  alarm: fakeAlarm(),
  alarms: fakeAlarms(10),
  meta: fakeMeta(),
  params: fakeParams(),
  widgetId: Faker.datatype.number(),
  alarmId: Faker.datatype.number(),
  csvData: Faker.datatype.string(),
  exportData: Faker.datatype.string(),
  allIds: Faker.datatype.array(5),
};

describe('Entities alarm module', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  beforeAll(() => {
    jest.mock('@/i18n', () => ({
      t: key => key,
    }));
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  afterAll(() => {
    jest.unmock('@/i18n');
  });

  it('Mutate state after commit FETCH_LIST', () => {
    const { params, widgetId } = mockData;
    const state = cloneDeep(initialState);

    const fetchList = mutations[types.FETCH_LIST];

    fetchList(state, { widgetId, params });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { pending: true, fetchingParams: params },
      },
    });
  });

  it('Mutate state after commit FETCH_LIST_COMPLETED', () => {
    const { meta, params, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { pending: true, fetchingParams: params },
      },
    };

    const fetchListCompleted = mutations[types.FETCH_LIST_COMPLETED];

    fetchListCompleted(state, { widgetId, meta });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { fetchingParams: params, meta, pending: false },
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

    const fetchListFailed = mutations[types.FETCH_LIST_FAILED];

    fetchListFailed(state, { widgetId });

    expect(state).toEqual({
      widgets: {
        [widgetId]: { pending: false, fetchingParams: params, meta },
      },
    });
  });

  it('Get alarms data by widget id. Getter: getListByWidgetId', () => {
    const { alarms, widgetId, allIds } = mockData;
    const getList = jest.fn(() => alarms);
    const rootGetters = {
      'entities/getList': getList,
    };
    const state = {
      widgets: {
        [widgetId]: { allIds },
      },
    };

    const data = getters.getListByWidgetId(state, getters, {}, rootGetters)(widgetId);

    expect(data).toEqual(alarms);
    expect(getList).toHaveBeenCalledWith(ENTITIES_TYPES.alarm, allIds);
  });

  it('Get alarm data by alarm id. Getter: getItem', () => {
    const { alarms } = mockData;
    const ids = Faker.datatype.array();
    const getList = jest.fn(() => alarms);
    const rootGetters = {
      'entities/getList': getList,
    };

    const data = getters.getList(initialState, getters, {}, rootGetters)(ids);

    expect(data).toEqual(alarms);
    expect(getList).toHaveBeenCalledWith(ENTITIES_TYPES.alarm, ids);
  });

  it('Get alarm data by alarm id. Getter: getItem', () => {
    const { alarmId, alarm } = mockData;
    const getItem = jest.fn(() => alarm);
    const rootGetters = {
      'entities/getItem': getItem,
    };

    const data = getters.getItem(initialState, getters, {}, rootGetters)(alarmId);

    expect(data).toEqual(alarm);
    expect(getItem).toHaveBeenCalledWith(ENTITIES_TYPES.alarm, alarmId);
  });

  it('Get meta data by widget id. Getter: getPendingByWidgetId', () => {
    const { meta, widgetId } = mockData;
    const state = {
      widgets: {
        [widgetId]: { meta },
      },
    };

    const data = getters.getMetaByWidgetId(state)(widgetId);

    expect(data).toEqual(meta);
  });

  it('Get meta data by widget id without data in state. Getter: getMetaByWidgetId', () => {
    const { widgetId } = mockData;

    const data = getters.getMetaByWidgetId(initialState)(widgetId);

    expect(data).toEqual({});
  });

  it('Get pending data by widget id. Getter: getPendingByWidgetId', () => {
    const { widgetId } = mockData;
    const pending = true;
    const state = {
      widgets: {
        [widgetId]: { pending },
      },
    };

    const data = getters.getPendingByWidgetId(state)(widgetId);

    expect(data).toEqual(pending);
  });

  it('Fetch alarms without saving in store and without params. Action: fetchListWithoutStore', async () => {
    axiosMockAdapter
      .onGet(API_ROUTES.alarms.list)
      .reply(200, mockData.alarmsResponse);

    const result = await actions.fetchListWithoutStore({}, {});

    expect(result).toEqual(mockData.alarmsResponse);
  });

  it('Fetch alarms without saving in store. Action: fetchListWithoutStore', async () => {
    const { params } = mockData;

    axiosMockAdapter
      .onGet(API_ROUTES.alarms.list, { params })
      .reply(200, mockData.alarmsResponse);

    const result = await actions.fetchListWithoutStore({}, { params });

    expect(result).toEqual(mockData.alarmsResponse);
  });

  it('Fetch alarm by id. Action: fetchItem', async () => {
    const dispatch = jest.fn();
    const { alarmId, params } = mockData;

    const data = await actions.fetchItem(
      { dispatch },
      { id: alarmId, params },
    );

    expect(data).toBeUndefined();
    expect(dispatch).toHaveBeenCalled();
  });
});
