import { mount, createVueInstance } from '@unit/utils/vue';
import ClickOutside from '@/services/click-outside';
import { ALARM_INTERVAL_FIELDS, QUICK_RANGES } from '@/constants';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

const localVue = createVueInstance();

const stubs = {
  'date-time-picker-text-field': true,
  'c-quick-date-interval-type-field': true,
};

const snapshotFactory = (options = {}) => mount(DateIntervalSelector, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

describe('date-interval-selector', () => {
  test('Renders `date-interval-selector` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last3Hour.start,
          tstop: QUICK_RANGES.last3Hour.stop,
          time_field: ALARM_INTERVAL_FIELDS.creationDate,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `date-interval-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last2Days.start,
          tstop: QUICK_RANGES.last2Days.stop,
          time_field: ALARM_INTERVAL_FIELDS.creationDate,
        },
        roundHours: true,
        tstartRules: {
          required: true,
        },
        tstopRules: {
          required: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
