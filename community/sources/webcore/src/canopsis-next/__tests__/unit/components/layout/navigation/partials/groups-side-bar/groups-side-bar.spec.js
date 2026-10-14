import { mockPopups } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createEntitiesModule,
  createMockedStoreModules,
  createModalsModule,
  createNavigationModule,
  createViewModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { CRUD_ACTIONS, MAX_LIMIT, USER_PERMISSIONS } from '@/constants';

import GroupsSideBar from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar.vue';

const vBtnStub = {
  template: '<button @click="$attrs.onClick?.($event)"><slot /></button>',
};

const stubs = {
  'app-logo': true,
  'logged-users-count': true,
  'app-version': true,
  'c-draggable-list-field': true,
  'groups-side-bar-group': true,
  'groups-side-bar-playlists': true,
  'groups-settings-button': {
    emits: ['toggleEditingMode'],
    template: '<div class="groups-settings-button-stub" @click="$emit(\'toggleEditingMode\')" />',
  },
  'v-btn': vBtnStub,
};

const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');
const selectOverlayButtons = wrapper => wrapper.findAll('.side-bar__overlay button');
const selectCancelButton = wrapper => selectOverlayButtons(wrapper).at(1);
const selectNavigationDrawer = wrapper => wrapper.find('v-navigation-drawer-stub');

describe('groups-side-bar', () => {
  const $popups = mockPopups();

  const { navigationModule, isEditingMode, toggleEditingMode } = createNavigationModule();
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

  const createStoreWithGroups = (availableGroups) => {
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.technical.view]: { actions: [CRUD_ACTIONS.read] },
    });
    groups.mockReturnValue(availableGroups);
    isEditingMode.mockReturnValue(true);

    return createMockedStoreModules([
      navigationModule,
      modalsModule,
      authModule,
      entitiesModule,
      viewModule,
    ]);
  };

  it('Side bar opened after trigger navigation', async () => {
    const wrapper = factory({ store });

    selectNavigationDrawer(wrapper).triggerCustomEvent('update:modelValue', true);

    expect(wrapper).toEmitInput(true);
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
        views: [
          {
            _id: Faker.datatype.string(),
          },
        ],
      },
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [
          {
            _id: Faker.datatype.string(),
          },
        ],
      },
    ];

    const wrapper = factory({
      store: createStoreWithGroups(availableGroups),
    });

    await flushPromises();

    fetchGroupsList.mockClear();

    const updatedGroups = [...availableGroups].reverse();

    wrapper.vm.mutatedGroups = updatedGroups.map(group => ({
      ...group,
      views: [...group.views],
    }));

    await wrapper.vm.submit();

    await flushPromises();

    expect(updateViewsPositions).toBeDispatchedWith({
      data: updatedGroups.map(group => ({
        _id: group._id,
        views: group.views.map(view => view._id),
      })),
    });
    expect(fetchGroupsList).toBeCalled();

    expect($popups.success).toHaveBeenCalledWith({
      text: 'layout.sideBar.ordering.popups.success',
    });
  });

  it('Error popup showed after trigger submit button with error', async () => {
    updateViewsPositions.mockRejectedValueOnce();

    const availableGroups = [
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{ _id: Faker.datatype.string() }],
      },
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [{ _id: Faker.datatype.string() }],
      },
    ];

    const wrapper = factory({
      store: createStoreWithGroups(availableGroups),
    });

    await flushPromises();

    fetchGroupsList.mockClear();

    const updatedGroups = [...availableGroups].reverse();

    wrapper.vm.mutatedGroups = updatedGroups.map(group => ({
      ...group,
      views: [...group.views],
    }));

    await wrapper.vm.submit();

    await flushPromises();

    expect(updateViewsPositions).toBeCalled();
    expect($popups.error).toHaveBeenCalledWith({ text: 'layout.sideBar.ordering.popups.error' });
  });

  it('Cancel groups order after trigger cancel button', async () => {
    const availableGroups = [
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [
          {
            _id: Faker.datatype.string(),
          },
        ],
      },
      {
        _id: Faker.datatype.string(),
        title: Faker.datatype.string(),
        views: [
          {
            _id: Faker.datatype.string(),
          },
        ],
      },
    ];

    const wrapper = factory({
      store: createStoreWithGroups(availableGroups),
    });

    await flushPromises();

    const updatedGroups = [...availableGroups].reverse();

    await selectDraggableField(wrapper).triggerCustomEvent('update:modelValue', updatedGroups);
    await selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect(wrapper.vm.isGroupsOrderChanged).toBe(false);
  });

  it('Editing mode changed after trigger settings button', async () => {
    const wrapper = factory({ store });

    await wrapper.vm.toggleNavigationEditingMode();

    expect(toggleEditingMode).toBeCalled();
  });
});
