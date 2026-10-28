import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { ref } from 'vue';

import ApiTokens from '@/components/other/api-token/api-tokens.vue';

const mocks = vi.hoisted(() => ({
  fetchApiTokensList: vi.fn(),
  removeApiToken: vi.fn(),
}));

vi.mock('@/hooks/store/modules/api-token', () => ({
  useApiToken: () => ({
    apiTokens: ref([]),
    apiTokensMeta: ref({ total_count: 0 }),
    apiTokensPending: ref(false),
    ...mocks,
  }),
}));

describe('api-tokens', () => {
  const factory = generateShallowRenderer(ApiTokens, {
    stubs: {
      'api-tokens-list': true,
    },
  });

  beforeEach(() => {
    mocks.fetchApiTokensList.mockReset().mockResolvedValue();
    mocks.removeApiToken.mockReset();
  });

  it('fetches only active public API tokens', async () => {
    factory();
    await flushPromises();

    expect(mocks.fetchApiTokensList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        active: true,
      },
    });
  });
});
