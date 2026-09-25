import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import {
  createEntityCategory,
  useEntityCategoryStore,
} from '@/stores/modules/entity-category';

const { requestGetMock, requestPostMock } = vi.hoisted(() => ({
  requestGetMock: vi.fn(),
  requestPostMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: requestGetMock,
    post: requestPostMock,
  },
}));

describe('entity category store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    requestGetMock.mockClear();
    requestPostMock.mockClear();
  });

  it('creates a category', () => {
    const data = { name: 'New category' };

    createEntityCategory({ data });

    expect(requestPostMock).toHaveBeenCalledWith(API_ROUTES.entityCategories, data);
  });

  it('serializes category ids without array brackets', () => {
    const params = {
      ids: ['first-category', 'second-category'],
      limit: 2,
    };

    useEntityCategoryStore().fetchCategoriesListWithoutStore({ params });

    expect(requestGetMock).toHaveBeenCalledWith(API_ROUTES.entityCategories, {
      params,
      paramsSerializer: { indexes: null },
    });
  });
});
