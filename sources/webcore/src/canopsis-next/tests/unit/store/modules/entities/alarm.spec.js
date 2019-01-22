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

describe('Alarm module', () => {
  const mockAdapter = new AxiosMockAdapter(request);

  beforeEach(() => {
    mockAdapter.reset();
  });

  it('fetchListWithoutStore', async (done) => {
    mockAdapter.onGet(API_ROUTES.alarmList).reply(200, mockData.alarmList);

    const result = await actions.fetchListWithoutStore({ dispatch: () => {} }, {});

    console.log(result);

    done();
  });

  it('fetchListWithoutStore', async (done) => {
    mockAdapter.onGet(API_ROUTES.alarmList).networkError();

    const result = await actions.fetchListWithoutStore({ dispatch: () => console.log('called dispatch') }, {});

    console.log(result);

    done();
  });
});
