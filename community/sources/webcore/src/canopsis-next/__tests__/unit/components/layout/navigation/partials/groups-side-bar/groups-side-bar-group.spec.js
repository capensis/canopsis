import { mockModals } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createMockedStoreModules,
  createNavigationModule,
  createActiveViewModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import GroupsSideBarGroup from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-group.vue';

const stubs = {
  'group-panel': true,
  'c-draggable-list-field': true,
  'groups-side-bar-group-view': true,
};

const selectGroupPanel = wrapper => wrapper.findComponent({ name: 'GroupPanel' });
const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');

describe('groups-side-bar-group', () => {
  const $modals = mockModals();

  const { navigationModule } = createNavigationModule();
  const { authModule } = createAuthModule();
  const { activeViewModule } = createActiveViewModule();
  const store = createMockedStoreModules([navigationModule, authModule, activeViewModule]);

  const group = {
    _id: 'group-1',
    views: [{ _id: 'group-view-1' }, { _id: 'group-view-2' }],
  };

  const factory = generateShallowRenderer(GroupsSideBarGroup, {
    stubs,
    mocks: { $modals, $route: { name: '', params: {} } },
  });

  test('Edit group modal showed after trigger group panel', async () => {
    const wrapper = factory({
      store,
      propsData: {
        group,
      },
    });

    await selectGroupPanel(wrapper).vm.$emit('change');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createGroup,
      config: {
        title: 'modals.group.edit.title',
        group,
        deletable: false,
        action: expect.any(Function),
        remove: expect.any(Function),
      },
    });
  });

  test('Change views order after trigger draggable list', async () => {
    const views = [
      { _id: Faker.datatype.string() },
      { _id: Faker.datatype.string() },
      { _id: Faker.datatype.string() },
    ];
    const groupWithViews = {
      _id: Faker.datatype.string(),
      views,
    };
    const wrapper = factory({
      store,
      propsData: {
        group: groupWithViews,
      },
    });

    const updatedViews = [...views].reverse();

    await selectDraggableField(wrapper).triggerCustomEvent('update:modelValue', updatedViews);

    expect(wrapper).toEmit('update:group', {
      ...groupWithViews,
      views: updatedViews,
    });
  });
});
