import { createPinia, setActivePinia } from 'pinia';

import { API_ROUTES } from '@/config';

import { useUserStore } from '@/stores/modules/user';

const { requestGetMock } = vi.hoisted(() => ({
  requestGetMock: vi.fn(),
}));

vi.mock('@/services/request', () => ({
  default: {
    get: requestGetMock,
  },
}));

describe('user store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    requestGetMock.mockReset();
  });

  it('fetches the user with the most active connections', () => {
    useUserStore().fetchTopByActiveConnectsWithoutStore();

    expect(requestGetMock).toHaveBeenCalledWith(API_ROUTES.topUserByActiveConnects);
  });
});
