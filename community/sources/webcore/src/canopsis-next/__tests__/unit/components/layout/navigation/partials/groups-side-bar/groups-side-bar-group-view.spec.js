import { mockModals } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createMockedStoreModules,
  createNavigationModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import GroupsSideBarGroupView from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-group-view.vue';

const stubs = {
  'group-view-panel': true,
  'router-link': {
    template: '<div><slot /></div>',
  },
};

const selectGroupViewPanel = wrapper => wrapper.findComponent({ name: 'GroupViewPanel' });

describe('groups-side-bar-group-view', () => {
  const $modals = mockModals();

  const { navigationModule } = createNavigationModule();
  const { authModule } = createAuthModule();
  const store = createMockedStoreModules([navigationModule, authModule]);

  const factory = generateShallowRenderer(GroupsSideBarGroupView, {
    stubs,
    mocks: { $modals },
  });

  it('Duplicate view modal showed after trigger duplicate event', async () => {
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

    await selectGroupViewPanel(wrapper).vm.$emit('duplicate');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createView,
      config: {
        title: 'modals.view.duplicate.title',
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

  it('Change view modal showed after trigger change event', async () => {
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

    await selectGroupViewPanel(wrapper).vm.$emit('change');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createView,
      config: {
        title: 'modals.view.edit.title',
        view,
        deletable: false,
        submittable: false,
        remove: expect.any(Function),
        action: expect.any(Function),
      },
    });
  });
});
