import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { API_ACCESS_TOKENS_TABS, MODALS } from '@/constants';

import ApiAccessTokensPage from '@/views/admin/api-access-tokens-page.vue';

const {
  createApiToken,
  createShareToken,
  fetchApiTokensListWithPreviousParams,
  fetchShareTokensListWithPreviousParams,
} = vi.hoisted(() => ({
  createApiToken: vi.fn(),
  createShareToken: vi.fn(),
  fetchApiTokensListWithPreviousParams: vi.fn().mockResolvedValue(),
  fetchShareTokensListWithPreviousParams: vi.fn().mockResolvedValue(),
}));

vi.mock('vue-router', () => ({
  useRouter: () => ({ resolve: vi.fn() }),
}));

vi.mock('@/hooks/auth', async () => {
  const { ref } = await vi.importActual('vue');

  return {
    useCRUDPermissions: () => ({
      hasCreateAccess: ref(true),
      hasReadAccess: ref(true),
    }),
  };
});

vi.mock('@/hooks/store/modules/api-token', () => ({
  useApiToken: () => ({
    createApiToken,
    fetchApiTokensListWithPreviousParams,
  }),
}));

vi.mock('@/hooks/store/modules/share-token', () => ({
  useShareToken: () => ({
    createShareToken,
    fetchShareTokensListWithPreviousParams,
  }),
}));

const stubs = {
  'c-page-header': true,
  'c-fab-btn': true,
  'api-tokens': true,
  'share-tokens': true,
  'v-card': { template: '<div><slot /></div>' },
  'v-card-text': { template: '<div><slot /></div>' },
  'v-tabs': { template: '<div><slot /></div>' },
  'v-tab': {
    name: 'VTab',
    props: ['value'],
    template: '<div><slot /></div>',
  },
  'v-window': { template: '<div><slot /></div>' },
  'v-window-item': { template: '<div><slot /></div>' },
};

describe('api-access-tokens-page', () => {
  const $modals = mockModals();
  const factory = generateShallowRenderer(ApiAccessTokensPage, { stubs });

  test('shows the public API tokens tab first and selects it by default', () => {
    const wrapper = factory({ mocks: { $modals } });
    const tabs = wrapper.findAllComponents({ name: 'VTab' });

    expect(wrapper.vm.activeTab).toBe(API_ACCESS_TOKENS_TABS.apiTokens);
    expect(tabs.map(tab => tab.props('value'))).toEqual([
      API_ACCESS_TOKENS_TABS.apiTokens,
      API_ACCESS_TOKENS_TABS.shareTokens,
    ]);
  });

  test('shows the one-time warning after creating a public API token', async () => {
    const token = { value: 'one-time-token' };
    createApiToken.mockResolvedValue(token);
    const wrapper = factory({ mocks: { $modals } });

    wrapper.vm.showCreateTokenModal();

    const createModal = $modals.show.mock.calls[0][0];

    expect(createModal.name).toBe(MODALS.createApiToken);

    await createModal.config.action({ name: 'token' });

    expect($modals.show).toHaveBeenLastCalledWith({
      name: MODALS.copyValue,
      config: {
        title: 'Public API token',
        label: 'Public API token',
        value: token.value,
        warningText:
          'Copy this token now. For security reasons, its full value will not be shown again.',
        copySuccessText: 'API token copied to clipboard',
      },
    });
  });
});
