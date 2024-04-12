import { generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimelineDays from '@/components/widgets/alarm/timeline/alarm-timeline-days.vue';

const stubs = {
  'alarm-timeline-steps': true,
};

describe('alarm-timeline-days', () => {
  const { steps: { data: steps } } = fakeAlarmDetails();
  const snapshotFactory = generateRenderer(AlarmTimelineDays, { stubs });

  test('Renders `alarm-timeline-days` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        days: groupAlarmSteps(steps),
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
