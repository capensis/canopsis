import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { computed, ref } from 'vue';

import Users from '@/components/other/users/users.vue';

vi.mock('@/hooks/auth', () => ({
  useAuth: () => ({
    currentUser: ref(),
    fetchCurrentUser: vi.fn(),
  }),
  useCRUDPermissions: () => ({
    hasUpdateAccess: computed(() => true),
    hasDeleteAccess: computed(() => true),
  }),
}));

vi.mock('@/hooks/store/modules/user', () => ({
  useUser: () => globalThis.TEST_USER_HOOK,
}));

const UsersListStub = {
  name: 'users-list',
  props: [
    'users',
    'options',
    'pending',
    'showTopUserByConnections',
  ],
  emits: [
    'edit',
    'remove',
    'refresh',
    'update:options',
    'update:showTopUserByConnections',
  ],
  template: '<div />',
};

describe('users', () => {
  const users = ref([{ _id: 'ordinary-user' }]);
  const usersMeta = ref({ total_count: 12 });
  const usersPending = ref(false);
  const fetchUsersList = vi.fn();
  const fetchTopUserByActiveConnectsWithoutStore = vi.fn();
  const updateUser = vi.fn();
  const removeUser = vi.fn();

  const factory = generateShallowRenderer(Users, {
    stubs: {
      'users-list': UsersListStub,
    },
  });

  beforeEach(() => {
    fetchUsersList.mockReset().mockResolvedValue();
    fetchTopUserByActiveConnectsWithoutStore.mockReset();
    updateUser.mockReset();
    removeUser.mockReset();

    globalThis.TEST_USER_HOOK = {
      users,
      usersMeta,
      usersPending,
      fetchUsersList,
      fetchTopUserByActiveConnectsWithoutStore,
      updateUser,
      removeUser,
    };
  });

  it('fetches the ordinary users list by default', async () => {
    factory();
    await flushPromises();

    expect(fetchUsersList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        with_flags: true,
      },
    });
    expect(fetchTopUserByActiveConnectsWithoutStore).not.toHaveBeenCalled();
  });

  it('resets search and shows only the top connected user when enabled', async () => {
    const topUser = { _id: 'top-user', active_connects: 7 };

    fetchTopUserByActiveConnectsWithoutStore.mockResolvedValue(topUser);

    const wrapper = factory();
    await flushPromises();

    wrapper.vm.updateListOptions({
      ...wrapper.vm.options,
      search: 'ordinary',
    });
    await flushPromises();

    fetchUsersList.mockClear();
    wrapper.vm.updateShowTopUserByConnections(true);
    await flushPromises();

    const list = wrapper.findComponent(UsersListStub);

    expect(fetchTopUserByActiveConnectsWithoutStore).toHaveBeenCalledTimes(1);
    expect(fetchUsersList).not.toHaveBeenCalled();
    expect(list.props('options').search).toBe('');
    expect(list.props('showTopUserByConnections')).toBe(true);
    expect(list.props('users')).toEqual([topUser]);

    wrapper.vm.updateListOptions({
      ...wrapper.vm.options,
      page: 2,
    });
    await flushPromises();

    expect(list.props('options').page).toBe(1);
    expect(fetchTopUserByActiveConnectsWithoutStore).toHaveBeenCalledTimes(1);
  });

  it('restores the ordinary users request when the filter is disabled', async () => {
    fetchTopUserByActiveConnectsWithoutStore.mockResolvedValue({ _id: 'top-user' });

    const wrapper = factory();
    await flushPromises();

    wrapper.vm.updateShowTopUserByConnections(true);
    await flushPromises();
    fetchUsersList.mockClear();

    wrapper.vm.updateShowTopUserByConnections(false);
    await flushPromises();

    expect(fetchUsersList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        with_flags: true,
      },
    });
    expect(wrapper.findComponent(UsersListStub).props('users')).toEqual(users.value);
  });
});
