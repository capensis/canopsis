import AxiosMockAdapter from 'axios-mock-adapter';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import alarmModule from '@/store/modules/entities/alarm';

import { fakeAlarms } from '@/unit/data/alarm';

const { actions } = alarmModule;

const mockData = {
  alarmList: fakeAlarms(10),
};

jest.mock('@/i18n', () => ({
  t: key => key,
}));

describe('Entities alarm module actions', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  afterAll(() => {
    jest.unmock('@/i18n');
  });

  it('fetchListWithoutStore without params', async () => {
    const dispatch = jest.fn();

    axiosMockAdapter
      .onGet(API_ROUTES.alarmList)
      .reply(200, mockData.alarmList);

    const result = await actions.fetchListWithoutStore({ dispatch }, {});

    expect(result).toEqual(mockData.alarmList);
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('fetchListWithoutStore with params', async () => {
    const dispatch = jest.fn();
    const payload = { params: { limit: 10 } };

    axiosMockAdapter
      .onGet(API_ROUTES.alarmList, payload)
      .reply(200, mockData.alarmList);

    const result = await actions.fetchListWithoutStore({ dispatch }, payload);

    expect(result).toEqual(mockData.alarmList);
    expect(dispatch).not.toHaveBeenCalled();
  });
});
