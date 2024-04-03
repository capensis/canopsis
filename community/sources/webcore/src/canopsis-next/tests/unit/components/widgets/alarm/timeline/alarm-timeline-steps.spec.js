import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import AlarmTimelineSteps from '@/components/widgets/alarm/timeline/alarm-timeline-steps.vue';

const stubs = {
  'alarm-timeline-step': true,
  'alarm-timeline-step-expand': true,
};

const selectAlarmTimelineSteps = wrapper => wrapper.findAll('alarm-timeline-step-stub');
const selectAlarmTimelineStepExpands = wrapper => wrapper.findAll('alarm-timeline-step-expand-stub');

describe('alarm-timeline-steps', () => {
  const { steps: { data: steps } } = fakeAlarmDetails();
  const factory = generateShallowRenderer(AlarmTimelineSteps, { stubs });
  const snapshotFactory = generateRenderer(AlarmTimelineSteps, { stubs });

  test('Toggles the expanded state when expand button is clicked', async () => {
    const wrapper = factory({
      propsData: {
        steps,
      },
    });

    const firstAlarmTimelineStep = selectAlarmTimelineSteps(wrapper).at(0);
    firstAlarmTimelineStep.triggerCustomEvent('expand', true);

    await flushPromises();

    const firstAlarmTimelineStepExpand = selectAlarmTimelineStepExpands(wrapper).at(0);
    expect(firstAlarmTimelineStepExpand.props('expanded')).toBeTruthy();
  });

  test('Toggles the expanded state when expand button is clicked for second step', async () => {
    const wrapper = factory({
      propsData: {
        steps,
      },
    });

    const secondAlarmTimelineStep = selectAlarmTimelineSteps(wrapper).at(1);
    secondAlarmTimelineStep.triggerCustomEvent('expand', true);

    await flushPromises();

    const firstAlarmTimelineStepExpand = selectAlarmTimelineStepExpands(wrapper).at(0);
    const secondAlarmTimelineStepExpand = selectAlarmTimelineStepExpands(wrapper).at(1);

    expect(firstAlarmTimelineStepExpand.props('expanded')).toBeFalsy();
    expect(secondAlarmTimelineStepExpand.props('expanded')).toBeTruthy();
  });

  test('Renders `alarm-timeline-steps` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
