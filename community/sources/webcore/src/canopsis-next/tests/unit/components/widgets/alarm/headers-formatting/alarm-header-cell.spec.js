import { generateRenderer } from '@unit/utils/vue';

import { ALARM_FIELDS } from '@/constants';

import AlarmHeaderCell from '@/components/widgets/alarm/headers-formatting/alarm-header-cell.vue';

const stubs = {
  'alarm-header-priority': true,
};

describe('alarm-header-cell', () => {
  const snapshotFactory = generateRenderer(AlarmHeaderCell, { stubs });

  it('Renders `alarm-header-cell` with priority header', () => {
    const wrapper = snapshotFactory({
      propsData: {
        header: {
          value: ALARM_FIELDS.impactState,
          text: 'Priority',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
