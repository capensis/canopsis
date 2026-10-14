import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { SIDE_BARS } from '@/constants';

import TheSidebar from '@/components/layout/sidebar/the-sidebar.vue';

const sidebarState = vi.hoisted(() => ({ sidebars: [] }));

vi.mock('@/stores/modules/sidebar', async () => {
  const { defineStore } = await vi.importActual('pinia');

  return {
    useSidebarStore: defineStore('sidebar', {
      state: () => sidebarState,
    }),
  };
});

const stubs = {
  'sidebar-base': {
    props: ['sidebar'],
    template: '<div class="sidebar-base">{{sidebar?.name}}</div>',
  },
};

describe('the-sidebar', () => {
  const factory = generateRenderer(TheSidebar, { stubs });

  it('Renders `the-sidebars` with type: alarmSettings', async () => {
    const sidebar = {
      id: 'test-sidebar-id',
      name: SIDE_BARS.alarmSettings,
      config: {},
      hidden: false,
    };

    sidebarState.sidebars = [sidebar];

    const wrapper = factory();

    await flushPromises();

    const sidebarBase = wrapper.find('.sidebar-base');

    expect(sidebarBase.text()).toEqual(sidebar.name);
  });
});
