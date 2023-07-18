import { generateRenderer } from '@unit/utils/vue';
import {
  createAuthModule, createEntitiesModule,
  createMockedStoreModules,
  createModalsModule,
  createNavigationModule,
  createViewGroupModule,
} from '@unit/utils/store';

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

describe('groups-side-bar', () => {
  const { viewGroupModule } = createViewGroupModule();
  const { navigationModule } = createNavigationModule();
  const { modalsModule } = createModalsModule();
  const { authModule } = createAuthModule();
  const { entitiesModule } = createEntitiesModule();
  const store = createMockedStoreModules([
    viewGroupModule,
    navigationModule,
    modalsModule,
    authModule,
    entitiesModule,
  ]);

  const snapshotFactory = generateRenderer(GroupsSideBar, { stubs });

  it('Renders `groups-side-bar` with closed navigation', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `groups-side-bar` with opened navigation', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        value: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
