import { createButtonStub } from '@unit/stubs/button';
import { mockModals, mockSidebar } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { SIDE_BARS } from '@/constants';

import SidebarWrapper from '@/plugins/sidebar/components/sidebar-wrapper.vue';

const stubs = {
  'v-btn': createButtonStub('v-btn'),
  'v-navigation-drawer': {
    props: ['value'],
    model: {
      event: 'input',
      prop: 'value',
    },
    template: `
      <div class="v-navigation-drawer">
        <input
          :checked="value"
          type="checkbox"
          @change="$emit('input', $event.target.checked)"
        />
        <slot />
      </div>
    `,
  },
};

describe('sidebar-wrapper', () => {
  const $modals = mockModals();
  const $sidebar = mockSidebar();

  const store = createMockedStoreModules([
    {
      name: $modals.moduleName,
      getters: {
        hasMaximizedModal: false,
      },
    },
  ]);

  const sidebar = {
    id: 'test-sidebar-id',
    name: SIDE_BARS.alarmSettings,
    config: {},
    hidden: false,
  };

  const factory = generateShallowRenderer(SidebarWrapper, {
    parentComponent: {
      provide: {
        $clickOutside: {
          call: vi.fn(),
        },
      },
    },

    stubs,
  });

  beforeEach(() => {
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((fn) => {
      fn(0);

      return 0;
    });
  });

  afterEach(() => {
    window.requestAnimationFrame.mockRestore();
  });

  it('Sidebar hidden trigger drawer', async () => {
    const wrapper = factory({
      propsData: {
        sidebar,
      },
      store,
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await flushPromises();

    const navigationDrawerInput = wrapper.find('.v-navigation-drawer input');

    navigationDrawerInput.setChecked(false);

    expect($sidebar.hide).toHaveBeenCalledTimes(1);
    expect($sidebar.hide).toHaveBeenCalledWith({ id: sidebar.id });
  });

  it('Sidebar hidden happened after click on the close button', async () => {
    const $clickOutside = {
      call: vi.fn(() => true),
    };

    const wrapper = factory({
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
      propsData: {
        sidebar,
      },
      store,
      mocks: {
        $modals,
        $sidebar,
      },
    });

    const closeButton = wrapper.findAll('button.v-btn').at(0);

    closeButton.trigger('click');

    expect($clickOutside.call).toHaveBeenCalledTimes(1);
    expect($sidebar.hide).toHaveBeenCalledTimes(1);
    expect($sidebar.hide).toHaveBeenCalledWith({ id: sidebar.id });
  });

  it('Sidebar hidden after click on the close button with close condition', async () => {
    const $clickOutside = {
      call: vi.fn(),
    };

    const wrapper = factory({
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
      propsData: {
        sidebar,
      },
      store,
      mocks: {
        $modals,
        $sidebar,
      },
    });

    const closeButton = wrapper.findAll('button.v-btn').at(0);

    closeButton.trigger('click');

    expect($clickOutside.call).toHaveBeenCalledTimes(1);
    expect($sidebar.hide).not.toHaveBeenCalled();
  });
});
