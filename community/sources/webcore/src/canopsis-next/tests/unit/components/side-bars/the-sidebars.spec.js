import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockRequestAnimationFrame } from '@unit/utils/mock-hooks';
import { SIDE_BARS } from '@/constants';

import TheSidebars from '@/plugins/sidebar/components/the-sidebar.vue';

const localVue = createVueInstance();

const snapshotStubs = {
  'alarm-settings': true,
  'context-settings': true,
  'service-weather-settings': true,
  'stats-histogram-settings': true,
  'stats-curves-settings': true,
  'stats-table-settings': true,
  'stats-calendar-settings': true,
  'stats-number-settings': true,
  'stats-pareto-settings': true,
  'text-settings': true,
  'counter-settings': true,
  'testing-weather-settings': true,
};

const snapshotFactory = (options = {}) => mount(TheSidebars, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('the-sidebars', () => {
  mockRequestAnimationFrame();

  it.each(Object.values(SIDE_BARS))('Renders `the-sidebars` with type: %s', async (type) => {
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
});
