import { cloneDeep, keyBy } from 'lodash';
import Vue from 'vue';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { fakeAlarm, fakeAlarms, fakeAlarmsResponse } from '@unit/data/alarm';
import { fakeMeta, fakeParams } from '@unit/data/request-data';

import { API_ROUTES } from '@/config';

import SetSeveralPlugin from '@/plugins/set-several';

import request from '@/services/request';

import alarmModule, { types } from '@/store/modules/entities/alarm';

import { mapIds } from '@/helpers/array';

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
      alarmsById: {},
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
    const { alarms, widgetId } = mockData;
    const state = {
      alarmsById: keyBy(alarms, '_id'),
      widgets: {
        [widgetId]: { allIds: mapIds(alarms) },
      },
    };

    const data = getters.getListByWidgetId(state, {
      ...getters,
      getList: getters.getList(state, {
        getItem: getters.getItem(state),
      }),
    })(widgetId);

    expect(data).toEqual(alarms);
  });

  it('Get alarms data by ids. Getter: getList', () => {
    const { alarms } = mockData;
    const ids = mapIds(alarms);
    const state = {
      ...initialState,
      alarmsById: keyBy(alarms, '_id'),
    };

    const data = getters.getList(state, {
      getItem: getters.getItem(state),
    })(ids);

    expect(data).toEqual(alarms);
  });

  it('Get alarm data by alarm id. Getter: getItem', () => {
    const { alarm } = mockData;

    const data = getters.getItem({
      ...initialState,
      alarmsById: {
        [alarm._id]: alarm,
      },
    })(alarm._id);

    expect(data).toEqual(alarm);
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
    const { params, alarm } = mockData;

    axiosMockAdapter
      .onGet(`${API_ROUTES.alarms.list}/${alarm._id}`)
      .reply(200, alarm);

    const commit = jest.fn();

    const data = await actions.fetchItem(
      { commit },
      { id: alarm._id, params },
    );

    expect(data).toEqual(alarm);
    expect(commit).toHaveBeenCalledWith(types.SET_ALARMS, { data: [alarm] });
  });
});
