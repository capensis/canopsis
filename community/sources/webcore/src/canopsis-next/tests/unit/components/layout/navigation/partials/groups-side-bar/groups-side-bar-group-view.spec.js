import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules, createNavigationModule } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import { CRUD_ACTIONS, MODALS } from '@/constants';

import GroupsSideBarGroupView from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-group-view.vue';

const stubs = {
  'group-view-panel': true,
  'router-link': true,
};

const selectGroupViewPanel = wrapper => wrapper.find('group-view-panel-stub');

describe('groups-side-bar-group-view', () => {
  const $modals = mockModals();

  const { navigationModule, isEditingMode } = createNavigationModule();
  const { authModule, currentUserPermissionsById } = createAuthModule();
  const store = createMockedStoreModules([navigationModule, authModule]);

  const factory = generateShallowRenderer(GroupsSideBarGroupView, {
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(GroupsSideBarGroupView, { stubs });

  it('Duplicate view modal showed after trigger duplicate event', () => {
    const view = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        view,
      },
      mocks: { $route: { params: { id: '' } } },
    });

    selectGroupViewPanel(wrapper).vm.$emit('duplicate');

    expect($modals.show).toBeCalledWith({
      name: MODALS.createView,
      config: {
        title: `Duplicate the view - ${view.title}`,
        duplicate: true,
        submittable: false,
        duplicableToAll: false,
        view: {
          ...view,
          name: '',
          title: '',
        },
        action: expect.any(Function),
      },
    });
  });

  it('Change view modal showed after trigger change event', () => {
    const view = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        view,
      },
      mocks: { $route: { params: { id: '' } } },
    });

    selectGroupViewPanel(wrapper).vm.$emit('change');

    expect($modals.show).toBeCalledWith({
      name: MODALS.createView,
      config: {
        title: 'Edit the view',
        view,
        deletable: false,
        submittable: false,
        remove: expect.any(Function),
        action: expect.any(Function),
      },
    });
  });

  it('Renders `groups-side-bar-group-view` with default data', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        view: { _id: 'view-1', title: 'View title' },
        isGroupsOrderChanged: true,
      },
      mocks: { $route: { params: { id: '' } } },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `groups-side-bar-group-view` with custom data', () => {
    const viewId = 'custom-view';

    isEditingMode.mockReturnValueOnce(true);
    currentUserPermissionsById.mockReturnValueOnce({
      [viewId]: { actions: [CRUD_ACTIONS.update] },
    });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([navigationModule, authModule]),
      propsData: {
        view: { _id: viewId, title: 'Custom view title' },
        isGroupsOrderChanged: true,
      },
      mocks: { $route: { params: { id: viewId } } },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
