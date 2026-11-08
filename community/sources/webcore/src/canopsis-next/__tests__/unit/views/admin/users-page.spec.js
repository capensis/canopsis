import { mockModals } from '@unit/utils/mock-hooks';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { ref } from 'vue';

import UsersList from '@/components/other/users/users-list.vue';

import UsersPage from '@/views/admin/users-page.vue';

const {
  fetchUsersList,
  fetchUsersListWithPreviousParams,
  fetchTopUserByActiveConnectsWithoutStore,
  createUser,
} = vi.hoisted(() => ({
  fetchUsersList: vi.fn(),
  fetchUsersListWithPreviousParams: vi.fn(),
  fetchTopUserByActiveConnectsWithoutStore: vi.fn(),
  createUser: vi.fn(),
}));

vi.mock('@/hooks/auth', () => ({
  useAuth: () => ({
    currentUser: ref(),
    fetchCurrentUser: vi.fn(),
  }),
  useCRUDPermissions: () => ({
    hasCreateAccess: ref(true),
    hasUpdateAccess: ref(true),
    hasDeleteAccess: ref(true),
  }),
}));

vi.mock('@/hooks/store/modules/user', () => ({
  useUser: () => ({
    users: ref([]),
    usersMeta: ref({ total_count: 0 }),
    usersPending: ref(false),
    fetchUsersList,
    fetchUsersListWithPreviousParams,
    fetchTopUserByActiveConnectsWithoutStore,
    createUser,
  }),
}));

const CPageStub = {
  name: 'c-page',
  emits: ['refresh', 'create'],
  template: '<div><slot /></div>',
};

const UsersListStub = {
  name: 'users-list',
  props: ['users', 'options'],
  emits: ['update:options', 'update:showTopUserByConnections'],
  template: '<div />',
};

describe('users-page', () => {
  const $modals = mockModals();
  const factory = generateRenderer(UsersPage, {
    stubs: {
      'c-page': CPageStub,
      'users-list': UsersListStub,
    },
  });

  beforeEach(() => {
    fetchUsersList.mockResolvedValue();
    fetchTopUserByActiveConnectsWithoutStore.mockResolvedValue({ _id: 'top-user' });
  });

  it('refreshes the top connected user when the toggle is enabled', async () => {
    const wrapper = factory();
    await flushPromises();
    const list = wrapper.findComponent(UsersList);

    list.vm.$emit('update:showTopUserByConnections', true);
    await flushPromises();
    fetchUsersList.mockClear();
    fetchTopUserByActiveConnectsWithoutStore.mockClear();

    const refreshedUser = { _id: 'new-top-user', active_connects: 12 };
    fetchTopUserByActiveConnectsWithoutStore.mockResolvedValueOnce(refreshedUser);

    wrapper.findComponent(CPageStub).vm.$emit('refresh');
    await flushPromises();

    expect(fetchTopUserByActiveConnectsWithoutStore).toHaveBeenCalledTimes(1);
    expect(fetchUsersList).not.toHaveBeenCalled();
    expect(fetchUsersListWithPreviousParams).not.toHaveBeenCalled();
    expect(list.props('users')).toEqual([refreshedUser]);
  });

  it('refreshes the ordinary users with the current table options when the toggle is disabled', async () => {
    const wrapper = factory();
    await flushPromises();
    const list = wrapper.findComponent(UsersList);

    list.vm.$emit('update:showTopUserByConnections', true);
    await flushPromises();
    list.vm.$emit('update:showTopUserByConnections', false);
    await flushPromises();
    list.vm.$emit('update:options', { ...list.props('options'), page: 2, search: 'alice' });
    await flushPromises();
    fetchUsersList.mockClear();
    fetchTopUserByActiveConnectsWithoutStore.mockClear();

    wrapper.findComponent(CPageStub).vm.$emit('refresh');
    await flushPromises();

    expect(fetchUsersList).toHaveBeenCalledExactlyOnceWith({
      params: { page: 2, limit: 10, search: 'alice', with_flags: true },
    });
    expect(fetchTopUserByActiveConnectsWithoutStore).not.toHaveBeenCalled();
  });

  it('refreshes the selected data source after creating a user', async () => {
    const wrapper = factory();
    await flushPromises();
    wrapper.findComponent(UsersList).vm.$emit('update:showTopUserByConnections', true);
    await flushPromises();
    fetchUsersList.mockClear();
    fetchTopUserByActiveConnectsWithoutStore.mockClear();

    wrapper.findComponent(CPageStub).vm.$emit('create');
    const data = { name: 'alice' };
    await $modals.show.mock.calls[0][0].config.action(data);

    expect(createUser).toHaveBeenCalledWith({ data });
    expect(fetchTopUserByActiveConnectsWithoutStore).toHaveBeenCalledTimes(1);
    expect(fetchUsersList).not.toHaveBeenCalled();
    expect(fetchUsersListWithPreviousParams).not.toHaveBeenCalled();
  });
});
