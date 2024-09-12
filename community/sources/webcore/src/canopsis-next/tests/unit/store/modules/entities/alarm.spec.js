import { keyBy } from 'lodash';
import Vue from 'vue';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { fakeAlarm, fakeAlarms, fakeAlarmsResponse } from '@unit/data/alarm';
import { fakeMeta, fakeParams } from '@unit/data/request-data';
import { API_ROUTES } from '@/config';

import SetSeveralPlugin from '@/plugins/set-several';

import request from '@/services/request';

import { getters as activeWidgetsGetters } from '@/store/modules/active-view/active-widgets';

import alarmModule, { types } from '@/store/modules/entities/alarm';

import { mapIds } from '@/helpers/array';

const { actions, state: initialState, getters } = alarmModule;

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

  it('Get alarms data by widget id. Getter: getListByWidgetId', () => {
    const { alarms, widgetId } = mockData;
    const state = {
      alarmsById: keyBy(alarms, '_id'),
    };

    const data = getters.getListByWidgetId(
      state,
      {
        ...getters,
        getList: getters.getList(state, {
          getItem: getters.getItem(state),
        }),
      },
      null,
      {
        [activeWidgetsGetters.GET_ALL_IDS_BY_WIDGET_ID]: () => mapIds(alarms),
      },
    )(widgetId);

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

  it('Get meta data by widget id. Getter: getMetaByWidgetId', () => {
    const { meta, widgetId } = mockData;

    const data = getters.getMetaByWidgetId(
      null,
      null,
      null,
      {
        [activeWidgetsGetters.GET_META_BY_WIDGET_ID]: () => meta,
      },
    )(widgetId);

    expect(data).toEqual(meta);
  });

  it('Get pending data by widget id. Getter: getPendingByWidgetId', () => {
    const { widgetId } = mockData;
    const pending = true;
    const data = getters.getPendingByWidgetId(
      null,
      null,
      null,
      {
        [activeWidgetsGetters.GET_PENDING_BY_WIDGET_ID]: () => pending,
      },
    )(widgetId);

    expect(data).toEqual(pending);
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
