import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules, createNavigationModule } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';

import GroupsSideBarGroup from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-group.vue';

const stubs = {
  'group-panel': true,
  'c-draggable-list-field': true,
  'groups-side-bar-group-view': true,
};

const selectGroupPanelNode = wrapper => wrapper.findRoot();
const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');

describe('groups-side-bar-group', () => {
  const $modals = mockModals();

  const { navigationModule } = createNavigationModule();
  const { authModule } = createAuthModule();
  const store = createMockedStoreModules([
    navigationModule,
    authModule,
  ]);

  const group = {
    views: [{ _id: 'group-view-1' }, { _id: 'group-view-2' }],
  };

  const factory = generateShallowRenderer(GroupsSideBarGroup, {
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(GroupsSideBarGroup, { stubs });

  it('Edit group modal showed after trigger group panel', () => {
    const wrapper = factory({
      store,
      propsData: {
        group,
      },
    });

    selectGroupPanelNode(wrapper).$emit('change');

    expect($modals.show).toBeCalledWith({
      name: MODALS.createGroup,
      config: {
        title: 'Edit group',
        group,
        deletable: false,
        action: expect.any(Function),
        remove: expect.any(Function),
      },
    });
  });

  it('Change views order after trigger draggable list', () => {
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

    selectDraggableField(wrapper).triggerCustomEvent('input', updatedViews);

    expect(wrapper).toEmit('update:group', {
      ...groupWithViews,
      views: updatedViews,
    });
  });

  it('Renders `groups-side-bar-group` with default data', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        group,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `groups-side-bar-group` with custom data', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        group,
        isGroupsOrderChanged: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `groups-side-bar-group` with empty groups', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        group: {
          ...group,
          views: [],
        },
        isGroupsOrderChanged: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
