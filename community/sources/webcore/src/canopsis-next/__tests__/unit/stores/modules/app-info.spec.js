import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import { i18n } from '@/i18n';

import { useAppInfoStore } from '@/stores/modules/app-info';
import { useI18nStore } from '@/stores/modules/i18n';

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
    i18n.global.locale.value = 'en';
  });

  afterEach(() => {
    i18n.global.locale.value = 'en';
  });

  it('applies the instance language when fetching app info without a callback', async () => {
    requestGetMock.mockResolvedValue({ language: 'fr' });

    await useAppInfoStore().fetchAppInfo();

    expect(i18n.global.locale.value).toBe('fr');
    expect(i18n.global.t('common.connect')).toBe('Connexion');
  });

  it('preserves the personal language when app info is refreshed', async () => {
    useI18nStore().setPersonalLocale('en');
    requestGetMock.mockResolvedValue({ language: 'fr' });

    await useAppInfoStore().fetchAppInfo();

    expect(i18n.global.locale.value).toBe('en');
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
