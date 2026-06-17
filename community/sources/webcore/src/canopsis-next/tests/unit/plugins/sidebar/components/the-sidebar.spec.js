import VueRouter from 'vue-router';

import { createVueInstance, flushPromises, generateRenderer } from '@unit/utils/vue';
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

const localVue = createVueInstance();

localVue.use(VueRouter);

const router = new VueRouter({
  mode: 'abstract',
  routes: [{ path: '/', component: { template: '<div />' } }],
});

describe('the-sidebar', () => {
  const $sidebar = mockSidebar();

  const snapshotFactory = generateRenderer(TheSidebar, { stubs: snapshotStubs, localVue, router });

  it('Renders `the-sidebars` with type: alarmSettings', async () => {
    const sidebar = {
      id: 'test-sidebar-id',
      name: SIDE_BARS.alarmSettings,
      config: {},
      hidden: false,
    };

    const wrapper = snapshotFactory({
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
