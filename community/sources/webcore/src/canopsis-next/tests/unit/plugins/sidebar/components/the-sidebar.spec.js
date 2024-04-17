import { flushPromises, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockSidebar } from '@unit/utils/mock-hooks';

import { SIDE_BARS } from '@/constants';

import TheSidebar from '@/plugins/sidebar/components/the-sidebar.vue';

const snapshotStubs = {
  'sidebar-base': {
    props: ['sidebar'],
    template: '<div class="sidebar-base">{{sidebar?.name}}</div>',
  },
};

describe('the-sidebar', () => {
  const $sidebar = mockSidebar();

  const snapshotFactory = generateRenderer(TheSidebar, { stubs: snapshotStubs });

  it('Renders `the-sidebars` with type: alarmSettings', async () => {
    const sidebar = {
      name: SIDE_BARS.alarmSettings,
      config: {},
      hidden: false,
    };

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          name: $sidebar.moduleName,
          getters: {
            sidebar,
          },
        },
      ]),
      mocks: {
        $sidebar,
      },
    });

    await flushPromises();

    const sidebarBase = wrapper.find('.sidebar-base');

    expect(sidebarBase.text()).toEqual(sidebar.name);
  });
});
