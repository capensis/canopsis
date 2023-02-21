import { mount, createVueInstance } from '@unit/utils/vue';

import { ALARM_FIELDS } from '@/constants';
import AlarmHeaderCell from '@/components/widgets/alarm/headers-formatting/alarm-header-cell.vue';

const localVue = createVueInstance();

const stubs = {
  'alarm-header-priority': true,
};

const snapshotFactory = (options = {}) => mount(AlarmHeaderCell, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-header-cell', () => {
  it('Renders `alarm-header-cell` with priority header', () => {
    const wrapper = snapshotFactory({
      propsData: {
        header: {
          value: ALARM_FIELDS.impactState,
          text: 'Priority',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-header-cell` with custom header', () => {
    const wrapper = snapshotFactory({
      propsData: {
        header: {
          value: 'custom-header',
          text: 'Custom header',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
