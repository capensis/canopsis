import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { COLOR_INDICATOR_TYPES } from '@/constants';

import AlarmColumnValue from '@/components/widgets/alarm/columns-formatting/alarm-column-value.vue';

const localVue = createVueInstance();

const stubs = {
  'color-indicator-wrapper': true,
  'alarm-column-cell': true,
  'c-runtime-template': true,
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
          colorIndicatorEnabled: true,
        },
        selectedTag: 'tag',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value` with custom template', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          name: 'alarm-name',
          entity: {
            name: 'entity-name',
          },
        },
        widget: {},
        column: {
          colorIndicator: COLOR_INDICATOR_TYPES.impactState,
          value: 'entity.name',
          template: '{{ value }} === {{ entity.name }} in the {{ alarm.name }}',
          colorIndicatorEnabled: true,
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
