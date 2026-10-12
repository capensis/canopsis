import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import { useAppInfoStore } from '@/stores/modules/app-info';

const { requestGetMock } = vi.hoisted(() => ({
  requestGetMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: requestGetMock,
  },
}));

describe('app info store', () => {
  beforeEach(() => {
    const pinia = createPinia();

    setActivePinia(pinia);

    requestGetMock.mockReset();
  });

  it('exposes the fetched UI date format unchanged', async () => {
    const dateFormat = 'Mon 02/01/2006 15:04:05';

    requestGetMock.mockResolvedValue({ date_format: dateFormat });

    const store = useAppInfoStore();

    await store.fetchAppInfo();

    expect(requestGetMock).toHaveBeenCalledWith(API_ROUTES.infos.app);
    expect(store.dateFormat).toBe(dateFormat);
  });

  it('leaves the UI date format empty when app-info omits date_format', async () => {
    requestGetMock.mockResolvedValue({});

    const store = useAppInfoStore();

    await store.fetchAppInfo();

    expect(store.dateFormat).toBeUndefined();
  });
});
