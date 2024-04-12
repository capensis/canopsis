import { generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import AlarmTimelineStepExpand from '@/components/widgets/alarm/timeline/alarm-timeline-step-expand.vue';

const stubs = {
  'alarm-timeline-step': true,
};

describe('alarm-timeline-step-expand', () => {
  const { steps: { data: steps } } = fakeAlarmDetails();
  const snapshotFactory = generateRenderer(AlarmTimelineStepExpand, { stubs });

  test('Renders `alarm-timeline-step-expand` with steps and without expanded correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: { steps },
        expanded: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-timeline-step-expand` with steps and expanded correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: { steps },
        expanded: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
