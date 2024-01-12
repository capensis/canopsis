import { generateRenderer } from '@unit/utils/vue';

import AlarmHeaderPriority from '@/components/widgets/alarm/headers-formatting/alarm-header-priority.vue';

const stubs = {
  'c-help-icon': true,
};

describe('alarm-header-priority', () => {
  const snapshotFactory = generateRenderer(AlarmHeaderPriority, { stubs });

  it('Renders `alarm-header-priority` without slot', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-header-priority` with slot', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: 'Default text slot',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
