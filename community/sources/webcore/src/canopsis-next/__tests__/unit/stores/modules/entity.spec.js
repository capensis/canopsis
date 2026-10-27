import { createPinia, setActivePinia } from 'pinia';

import { useEntityStore } from '@/stores/modules/entity';

vi.unmock('@/stores/modules/entity');

const { requestGetMock } = vi.hoisted(() => ({
  requestGetMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: requestGetMock,
  },
  isRequestCanceled: () => false,
  useRequestCancelling: callback => callback({ token: 'cancel-token' }),
}));

describe('entity store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    requestGetMock.mockReset();
  });

  it('replaces a cached entity with its refreshed value', async () => {
    const entityId = 'entity-id';
    const widgetId = 'widget-id';
    const store = useEntityStore();

    requestGetMock.mockResolvedValueOnce({
      data: [{
        _id: entityId,
        name: 'Old entity name',
        pause: { timestamp: 123 },
        obsolete_field: true,
      }],
      meta: { total_count: 1 },
    });

    await store.fetchList({ widgetId, params: {} });

    expect(store.getListByWidgetId(widgetId)[0].pause).toEqual({ timestamp: 123 });

    requestGetMock.mockResolvedValueOnce({
      data: [{ _id: entityId, name: 'New entity name' }],
      meta: { total_count: 1 },
    });

    await store.fetchListWithPreviousParams();

    expect(store.getListByWidgetId(widgetId)[0]).toEqual({
      _id: entityId,
      name: 'New entity name',
    });
  });
});
