import Faker from 'faker';
import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import {
  createAuthModule,
  createEntitiesModule,
  createMockedStoreModules,
  createModalsModule,
  createNavigationModule,
  createViewModule,
} from '@unit/utils/store';
import { CRUD_ACTIONS, MAX_LIMIT, USERS_PERMISSIONS } from '@/constants';
import { mockPopups } from '@unit/utils/mock-hooks';

import GroupsSideBar from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar.vue';

const stubs = {
  'app-logo': true,
  'logged-users-count': true,
  'app-version': true,
  'c-draggable-list-field': true,
  'groups-side-bar-group': true,
  'groups-side-bar-playlists': true,
  'groups-settings-button': true,
};

const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');
const selectButtons = wrapper => wrapper.findAll('v-btn-stub');
const selectGroupsSettingsButton = wrapper => wrapper.find('groups-settings-button-stub');
const selectNavigationDrawer = wrapper => wrapper.find('v-navigation-drawer-stub');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(0);
const selectCancelButton = wrapper => selectButtons(wrapper).at(1);

describe('groups-side-bar', () => {
  const $popups = mockPopups();

  const { navigationModule, toggleEditingMode } = createNavigationModule();
  const { modalsModule } = createModalsModule();
  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { entitiesModule } = createEntitiesModule();
  const { viewModule, updateViewsPositions, groups, fetchGroupsList } = createViewModule();
  const store = createMockedStoreModules([
    navigationModule,
    modalsModule,
    authModule,
    entitiesModule,
    viewModule,
  ]);

  const factory = generateShallowRenderer(GroupsSideBar, {
    stubs,
    mocks: { $popups },
  });
  const snapshotFactory = generateRenderer(GroupsSideBar, { stubs });

  it('Side bar opened after trigger navigation', async () => {
    const wrapper = factory({ store });

    selectNavigationDrawer(wrapper).vm.$emit('input', true);

    expect(wrapper).toEmit('input', true);
  });

  it('Groups fetched after mount', async () => {
    factory({ store });

    await flushPromises();

    expect(fetchGroupsList).toBeDispatchedWith({
      params: {
        limit: MAX_LIMIT,
        page: 1,
        with_views: true,
        with_tabs: true,
        with_widgets: true,
        with_flags: true,
        with_private: true,
      },
    });
  });

  it('Change groups order after trigger submit button', async () => {
    const availableGroups = [
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{
          _id: Faker.datatype.string(),
        }],
      },
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{
          _id: Faker.datatype.string(),
        }],
      },
    ];

    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.view]: { actions: [CRUD_ACTIONS.read] },
    });
    groups.mockReturnValueOnce(availableGroups);

    const wrapper = factory({
      store: createMockedStoreModules([
        navigationModule,
        modalsModule,
        authModule,
        entitiesModule,
        viewModule,
      ]),
    });

    await flushPromises();

    fetchGroupsList.mockClear();

    const updatedGroups = [...availableGroups].reverse();

    selectDraggableField(wrapper).vm.$emit('input', updatedGroups);
    selectSubmitButton(wrapper).vm.$emit('click');

    await flushPromises();

    expect(updateViewsPositions).toBeDispatchedWith({
      data: updatedGroups.map(group => ({
        _id: group._id,
        views: group.views.map(view => view._id),
      })),
    });
    expect(fetchGroupsList).toBeCalled();

    expect($popups.success).toBeCalledWith({ text: 'The groups was reordered' });
  });

  it('Error popup showed after trigger submit button with error', async () => {
    updateViewsPositions.mockRejectedValueOnce();
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.view]: { actions: [CRUD_ACTIONS.read] },
    });

    const wrapper = factory({
      store: createMockedStoreModules([
        navigationModule,
        modalsModule,
        authModule,
        entitiesModule,
        viewModule,
      ]),
    });

    await flushPromises();

    fetchGroupsList.mockClear();

    selectSubmitButton(wrapper).vm.$emit('click');

    await flushPromises();

    expect(updateViewsPositions).toBeCalled();
    expect($popups.error).toBeCalledWith({ text: 'Several groups wasn\'t reordered' });
  });

  it('Cancel groups order after trigger cancel button', async () => {
    const availableGroups = [
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{
          _id: Faker.datatype.string(),
        }],
      },
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{
          _id: Faker.datatype.string(),
        }],
      },
    ];

    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.view]: { actions: [CRUD_ACTIONS.read] },
    });
    groups.mockReturnValueOnce(availableGroups);

    const wrapper = factory({
      store: createMockedStoreModules([
        navigationModule,
        modalsModule,
        authModule,
        entitiesModule,
        viewModule,
      ]),
    });

    const updatedGroups = [...availableGroups].reverse();

    selectDraggableField(wrapper).vm.$emit('input', updatedGroups);
    selectCancelButton(wrapper).vm.$emit('click');

    await flushPromises();

    expect(selectSubmitButton(wrapper).isVisible()).toBe(false);
  });

  it('Editing mode changed after trigger settings button', async () => {
    const wrapper = factory({ store });

    selectGroupsSettingsButton(wrapper).vm.$emit('toggleEditingMode');

    expect(toggleEditingMode).toBeCalled();
  });

  it('Renders `groups-side-bar` with default data', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `groups-side-bar` with custom data', () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.view]: { actions: [CRUD_ACTIONS.read] },
    });
    groups.mockReturnValueOnce([
      { title: 'Group 1', views: [{}] },
      { title: 'Group 2', views: [{}] },
    ]);

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        navigationModule,
        modalsModule,
        authModule,
        entitiesModule,
        viewModule,
      ]),
      propsData: {
        value: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
