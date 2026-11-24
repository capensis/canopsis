import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import { useEntityInfoPropertyStore } from '@/stores/modules/entity-info-property';

vi.unmock('@/stores/modules/entity-info-property');

const { requestGetMock } = vi.hoisted(() => ({
  requestGetMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: requestGetMock,
  },
}));

describe('entity info property store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    requestGetMock.mockReset();
  });

  it('shares an active request and allows fetching again after completion', async () => {
    const store = useEntityInfoPropertyStore();
    const response = { data: [{ _id: 'property', alias: 'Alias' }], meta: { total_count: 1 } };
    const { promise, resolve } = Promise.withResolvers();

    requestGetMock.mockReturnValueOnce(promise).mockResolvedValue(response);

    const first = store.fetchList({ params: { paginate: false } });
    const second = store.fetchList({ params: { paginate: false } });
    const third = store.fetchListWithPreviousParams();

    expect(requestGetMock).toHaveBeenCalledExactlyOnceWith(API_ROUTES.entityInfosProperties, {
      params: { paginate: false },
    });
    expect(store.pending).toBe(true);

    resolve(response);

    expect(await Promise.all([first, second, third])).toEqual([response, response, response]);
    expect(store.items).toEqual(response.data);
    expect(store.itemsWithAlias).toEqual(response.data);
    expect(store.meta).toEqual(response.meta);
    expect(store.pending).toBe(false);

    await store.fetchList({ params: { paginate: false } });

    expect(requestGetMock).toHaveBeenCalledTimes(2);
  });

  it('matches parameters regardless of object key order', async () => {
    const store = useEntityInfoPropertyStore();
    const { promise, resolve } = Promise.withResolvers();

    requestGetMock.mockReturnValue(promise);

    const first = store.fetchList({ params: { paginate: false, type: 1 } });
    const second = store.fetchList({ params: { type: 1, paginate: false } });

    expect(requestGetMock).toHaveBeenCalledTimes(1);

    resolve({ data: [] });
    await Promise.all([first, second]);
  });

  it('keeps different requests separate and stays pending until both finish', async () => {
    const store = useEntityInfoPropertyStore();
    const firstRequest = Promise.withResolvers();
    const secondRequest = Promise.withResolvers();

    requestGetMock.mockReturnValueOnce(firstRequest.promise).mockReturnValueOnce(secondRequest.promise);

    const first = store.fetchList({ params: { paginate: false } });
    const second = store.fetchList({ params: { paginate: false, type: 1 } });
    const repeatedFirst = store.fetchList({ params: { paginate: false } });

    expect(requestGetMock).toHaveBeenCalledTimes(2);

    firstRequest.resolve({ data: [] });
    await Promise.all([first, repeatedFirst]);

    expect(store.pending).toBe(true);

    secondRequest.resolve({ data: [] });
    await second;

    expect(store.pending).toBe(false);
  });

  it('allows retrying after a shared request fails', async () => {
    const store = useEntityInfoPropertyStore();
    const error = new Error('Request failed');
    const { promise, reject } = Promise.withResolvers();
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

    requestGetMock.mockReturnValueOnce(promise).mockResolvedValue({ data: [] });

    const first = store.fetchList({ params: { paginate: false } });
    const second = store.fetchList({ params: { paginate: false } });
    const results = Promise.allSettled([first, second]);

    reject(error);

    expect(await results).toEqual([
      { status: 'rejected', reason: error },
      { status: 'rejected', reason: error },
    ]);
    expect(store.pending).toBe(false);

    await store.fetchList({ params: { paginate: false } });

    expect(requestGetMock).toHaveBeenCalledTimes(2);
    consoleError.mockRestore();
  });

  it('does not share requests between store instances', async () => {
    const firstStore = useEntityInfoPropertyStore();
    const secondStore = useEntityInfoPropertyStore(createPinia());

    requestGetMock.mockResolvedValue({ data: [] });

    await Promise.all([
      firstStore.fetchList({ params: { paginate: false } }),
      secondStore.fetchList({ params: { paginate: false } }),
    ]);

    expect(requestGetMock).toHaveBeenCalledTimes(2);
  });
});
