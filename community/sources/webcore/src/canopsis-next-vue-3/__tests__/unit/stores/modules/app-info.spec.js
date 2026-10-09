import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';
import {
  DEFAULT_USER_INTERFACE_DATE_FORMAT,
  setUserInterfaceDateFormat,
} from '@/constants';

import { useAppInfoStore } from '@/stores/modules/app-info';

import { convertDateToString } from '@/helpers/date/date';

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
    setActivePinia(createPinia());
    requestGetMock.mockReset();
  });

  afterEach(() => setUserInterfaceDateFormat(DEFAULT_USER_INTERFACE_DATE_FORMAT));

  it('applies the fetched UI date format unchanged', async () => {
    const dateFormat = 'Mon 02/01/2006 15:04:05';
    const date = new Date(2021, 8, 20, 21, 15, 40);

    requestGetMock.mockResolvedValue({ date_format: dateFormat });

    const store = useAppInfoStore();

    await store.fetchAppInfo();

    expect(requestGetMock).toHaveBeenCalledWith(API_ROUTES.infos.app);
    expect(store.dateFormat).toBe(dateFormat);
    expect(convertDateToString(date)).toBe('Mon 20/09/2021 21:15:40');
  });

  it('uses the default UI format when app-info omits date_format', async () => {
    setUserInterfaceDateFormat('Mon 02/01/2006 15:04:05');
    requestGetMock.mockResolvedValue({});

    await useAppInfoStore().fetchAppInfo();

    expect(convertDateToString(new Date(2021, 8, 20, 21, 15, 40))).toBe(
      '20/09/2021 21:15:40',
    );
  });
});
