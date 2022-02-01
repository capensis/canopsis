import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { mockRequestAnimationFrame } from '@unit/utils/mock-hooks';
import { SIDE_BARS } from '@/constants';

import SideBarWrapper from '@/components/sidebars/side-bar-wrapper.vue';

const localVue = createVueInstance();

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

const factory = (options = {}) => shallowMount(SideBarWrapper, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(SideBarWrapper, {
  localVue,

  ...options,
});

describe('side-bar-wrapper', () => {
  mockRequestAnimationFrame();

  it('Modal hidden trigger drawer', async () => {
    const hideSideBar = jest.fn();
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {},
            hidden: false,
          },
          actions: {
            hide: hideSideBar,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    await flushPromises();

    const navigationDrawerInput = wrapper.find('.v-navigation-drawer input');

    navigationDrawerInput.setChecked(false);

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Modal hidden after click on the close button', async () => {
    const hideSideBar = jest.fn();
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {},
            hidden: false,
          },
          actions: {
            hide: hideSideBar,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    const closeButton = wrapper.find('button.v-btn');

    closeButton.trigger('click');

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Modal hidden after click on the close button with close condition', async () => {
    const hideSideBar = jest.fn();
    const clickOutsideCondition = jest.fn(() => false);
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {},
            hidden: false,
          },
          actions: {
            hide: hideSideBar,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    const { $clickOutside } = wrapper.vm;

    $clickOutside.register(clickOutsideCondition);

    const closeButton = wrapper.find('button.v-btn');

    closeButton.trigger('click');

    await flushPromises();

    expect(clickOutsideCondition).toHaveBeenCalledTimes(1);
    expect(hideSideBar).not.toHaveBeenCalled();
  });

  it.each(Object.values(SIDE_BARS))('Renders `side-bar-wrapper` with type: %s', async (type) => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: type,
            config: {},
            hidden: false,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `side-bar-wrapper` with type custom title', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {
              sideBarTitle: 'Custom title',
            },
            hidden: false,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `side-bar-wrapper` with default slot', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {},
            hidden: false,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `side-bar-wrapper` with navigation drawer props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        navigationDrawerProps: {
          right: true,
          fixed: true,
          temporary: true,
          width: 400,
        },
      },
      store: createMockedStoreModules([
        {
          name: 'sideBar',
          getters: {
            name: SIDE_BARS.alarmSettings,
            config: {},
            hidden: false,
          },
        },
        {
          name: 'modals',
          getters: {
            hasMaximizedModal: false,
          },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
