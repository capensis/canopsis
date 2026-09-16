import { mockSidebar } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { SIDE_BARS } from '@/constants';

import TheSidebar from '@/plugins/sidebar/components/the-sidebar.vue';

const stubs = {
  'sidebar-base': {
    props: ['sidebar'],
    template: '<div class="sidebar-base">{{sidebar?.name}}</div>',
  },
};

describe('the-sidebar', () => {
  const $sidebar = mockSidebar();

  const factory = generateRenderer(TheSidebar, { stubs });

  it('Renders `the-sidebars` with type: alarmSettings', async () => {
    const sidebar = {
      id: 'test-sidebar-id',
      name: SIDE_BARS.alarmSettings,
      config: {},
      hidden: false,
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: $sidebar.moduleName,
          getters: {
            sidebars: [sidebar],
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
