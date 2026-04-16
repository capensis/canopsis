import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { mockModals, mockSidebar } from '@unit/utils/mock-hooks';

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
          @change="$listeners.input($event.target.checked)"
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
          call: jest.fn(),
        },
      },
    },

    stubs,
  });
  const snapshotFactory = generateRenderer(SidebarWrapper, {
    parentComponent: {
      provide: {
        $clickOutside: {
          call: jest.fn(),
        },
      },
    },
  });

  beforeEach(() => {
    jest.spyOn(window, 'requestAnimationFrame').mockImplementation((fn) => {
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
      call: jest.fn(() => true),
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
      call: jest.fn(),
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

  it('Sidebar minimize after click when config is minimizable', async () => {
    const wrapper = factory({
      propsData: {
        sidebar: {
          id: 'test-sidebar-id',
          name: SIDE_BARS.alarmSettings,
          config: {
            minimizable: true,
            title: 'Test title',
          },
          hidden: false,
        },
      },
      store,
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await flushPromises();

    wrapper.findAll('button.v-btn').at(0).trigger('click');

    expect($sidebar.minimize).toHaveBeenCalledTimes(1);
    expect($sidebar.minimize).toHaveBeenCalledWith({ id: 'test-sidebar-id' });
    expect($sidebar.hide).not.toHaveBeenCalled();
  });

  it('Sidebar maximize after click when minimized and minimizable', async () => {
    const wrapper = factory({
      propsData: {
        sidebar: {
          id: 'test-sidebar-id',
          name: SIDE_BARS.alarmSettings,
          config: {
            minimizable: true,
            title: 'Test title',
          },
          hidden: false,
          minimized: true,
        },
      },
      store,
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await flushPromises();

    wrapper.findAll('button.v-btn').at(1).trigger('click');

    expect($sidebar.maximize).toHaveBeenCalledTimes(1);
    expect($sidebar.maximize).toHaveBeenCalledWith({ id: 'test-sidebar-id' });
  });

  it.each(Object.values(SIDE_BARS))('Renders `sidebar-wrapper` with type: %s', async (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        sidebar: {
          id: 'test-sidebar-id',
          name: type,
          config: {},
          hidden: false,
        },
      },
      store: createMockedStoreModules([
        {
          name: $modals.moduleName,
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `sidebar-wrapper` with default slot', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        sidebar,
      },
      store,
      slots: {
        default: '<div class="default-slot" />',
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `sidebar-wrapper` without sidebar name', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        sidebar: {
          id: 'test-sidebar-id',
          config: {},
          hidden: false,
        },
      },
      store,
      slots: {
        default: '<div class="default-slot" />',
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
