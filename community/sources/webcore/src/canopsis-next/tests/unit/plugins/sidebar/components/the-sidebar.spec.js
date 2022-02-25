import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockSidebar } from '@unit/utils/mock-hooks';

import { SIDE_BARS } from '@/constants';

import TheSidebar from '@/plugins/sidebar/components/the-sidebar.vue';

const localVue = createVueInstance();

const stubs = {
  'sidebar-base': {
    props: ['sidebar'],
    template: '<div class="sidebar-base">{{sidebar?.name}}</div>',
  },
};

const snapshotFactory = (options = {}) => mount(TheSidebar, {
  localVue,
  stubs,

  ...options,
});

describe('the-sidebar', () => {
  const $sidebar = mockSidebar();

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
