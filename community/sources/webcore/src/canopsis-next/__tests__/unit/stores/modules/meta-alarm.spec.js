import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import { useMetaAlarmStore } from '@/stores/modules/meta-alarm';

vi.unmock('@/stores/modules/meta-alarm');

const { requestPutMock } = vi.hoisted(() => ({
  requestPutMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: requestPutMock,
  },
}));

describe('meta alarm store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    requestPutMock.mockClear();
  });

  it('adds alarms to an existing meta alarm', () => {
    const id = 'meta-alarm-id';
    const data = { alarms: ['alarm-id'] };

    useMetaAlarmStore().addAlarms({ id, data });

    expect(requestPutMock).toHaveBeenCalledWith(`${API_ROUTES.metaAlarm}/${id}/add`, data);
  });

  it('removes alarms from a meta alarm', () => {
    const id = 'meta-alarm-id';
    const data = { alarms: ['alarm-id'] };

    useMetaAlarmStore().removeAlarms({ id, data });

    expect(requestPutMock).toHaveBeenCalledWith(`${API_ROUTES.metaAlarm}/${id}/remove`, data);
  });
});
