import { mount, createVueInstance } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import AlarmColumnValue from '@/components/widgets/alarm/columns-formatting/alarm-column-value.vue';

const localVue = createVueInstance();

const stubs = {
  'color-indicator-wrapper': true,
  'alarm-column-cell': true,
};

const snapshotFactory = (options = {}) => mount(AlarmColumnValue, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-column-value', () => {
  it('Renders `alarm-column-value` with required props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
        },
        widget: {},
        column: {
          colorIndicator: COLOR_INDICATOR_TYPES.state,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
        },
        widget: {},
        column: {
          colorIndicator: COLOR_INDICATOR_TYPES.impactState,
        },
        columnsFilters: [{}, {}],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
