import { mount, createVueInstance } from '@unit/utils/vue';

import { EVENT_ENTITY_TYPES } from '@/constants';
import AlarmColumnValueState from '@/components/widgets/alarm/columns-formatting/alarm-column-value-state.vue';

const localVue = createVueInstance();

const stubs = {
  'c-alarm-chip': true,
};

const snapshotFactory = (options = {}) => mount(AlarmColumnValueState, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-column-value-state', () => {
  it('Renders `alarm-column-value-state` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-state` with alarm state', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          v: {
            events_count: 'Events count',
            state: {
              val: 'custom-state-val',
              _t: EVENT_ENTITY_TYPES.changeState,
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-state` with custom propertyKey', () => {
    const wrapper = snapshotFactory({
      propsData: {
        propertyKey: 'customPropertyKey',
        alarm: {
          customPropertyKey: 'custom-property-value',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
