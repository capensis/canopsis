import { generateShallowRenderer } from '@unit/utils/vue';

import UsersList from '@/components/other/users/users-list.vue';

const CAdvancedDataTableStub = {
  name: 'c-advanced-data-table',
  props: ['searchDisabled'],
  template: '<div class="advanced-data-table"><slot name="toolbar" /></div>',
};

const VSwitchStub = {
  name: 'v-switch',
  props: ['modelValue', 'label'],
  emits: ['update:modelValue'],
  template: '<button class="top-user-switch" @click="$emit(\'update:modelValue\', !modelValue)">{{ label }}</button>',
};

const stubs = {
  'c-advanced-data-table': CAdvancedDataTableStub,
  'c-table-mass-actions-panel': true,
  'c-enabled': true,
  'c-icon': true,
  'c-action-btn': true,
  'v-switch': VSwitchStub,
};

describe('users-list', () => {
  const factory = generateShallowRenderer(UsersList, { stubs });

  const createWrapper = showTopUserByConnections => factory({
    props: {
      options: {
        page: 1,
        itemsPerPage: 10,
        search: '',
      },
      showTopUserByConnections,
    },
  });

  it('disables search while the top-user filter is enabled', () => {
    const wrapper = createWrapper(true);

    expect(wrapper.findComponent(CAdvancedDataTableStub).props('searchDisabled')).toBe(true);
    expect(wrapper.find('.top-user-switch').text()).toBe('user.showTopUserByConnections');
  });

  it('emits the updated top-user filter value', async () => {
    const wrapper = createWrapper(false);

    await wrapper.find('.top-user-switch').trigger('click');

    expect(wrapper).toEmit('update:showTopUserByConnections', true);
  });
});
