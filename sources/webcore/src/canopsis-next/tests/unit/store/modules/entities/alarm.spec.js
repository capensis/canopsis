import AxiosMockAdapter from 'axios-mock-adapter';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import alarmModule from '@/store/modules/entities/alarm';

const { actions } = alarmModule;

const mockData = {
  alarmList: {
    data: [
      {
        alarms: [
          { title: 'Something' },
        ],
        first: 1,
        last: 10,
        total: 290,
        truncated: false,
      },
    ],
  },
};

jest.mock('@/i18n', () => ({
  t: key => key,
}));

describe('Alarm module', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  afterAll(() => {
    jest.unmock('@/i18n');
  });

  it('fetchListWithoutStore without params', async () => {
    const dispatch = jest.fn();
    const payload = {};

    axiosMockAdapter.onGet(API_ROUTES.alarmList).reply(200, mockData.alarmList);

    const result = await actions.fetchListWithoutStore({ dispatch }, payload);

    expect(result).toEqual(mockData.alarmList.data[0]);
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('fetchListWithoutStore with params', async () => {
    const dispatch = jest.fn();
    const payload = { params: { limit: 5, skip: 0 } };

    axiosMockAdapter.onGet(API_ROUTES.alarmList, payload.params).reply(200, mockData.alarmList);

    const result = await actions.fetchListWithoutStore({ dispatch }, payload);

    expect(result).toEqual(mockData.alarmList.data[0]);
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('fetchListWithoutStore with params', async () => {
    const dispatch = jest.fn();
    const payload = {};

    axiosMockAdapter.onGet(API_ROUTES.alarmList).networkError();

    const result = await actions.fetchListWithoutStore({ dispatch }, payload);

    expect(result).toEqual({ alarms: [], total: 0 });
    expect(dispatch).toHaveBeenCalledTimes(1);
    expect(dispatch).toHaveBeenCalledWith('popup/add', { type: 'error', text: 'errors.default' }, { root: true });
  });
});
