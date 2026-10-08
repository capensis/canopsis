import { fakeAlarmDetails } from '@unit/data/alarm';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import AlarmTimelineSteps from '@/components/widgets/alarm/timeline/alarm-timeline-steps.vue';

const alarmTimelineStepStub = {
  name: 'AlarmTimelineStep',
  emits: ['expand'],
  template: '<div />',
};
const alarmTimelineStepExpandStub = {
  name: 'AlarmTimelineStepExpand',
  props: ['expanded'],
  template: '<div />',
};

const stubs = {
  'alarm-timeline-step': alarmTimelineStepStub,
  'alarm-timeline-step-expand': alarmTimelineStepExpandStub,
};

const selectAlarmTimelineSteps = wrapper => wrapper.findAllComponents(alarmTimelineStepStub);
const selectAlarmTimelineStepExpands = wrapper =>
  wrapper.findAllComponents(alarmTimelineStepExpandStub);

describe('alarm-timeline-steps', () => {
  const {
    steps: { data: steps },
  } = fakeAlarmDetails();
  const factory = generateShallowRenderer(AlarmTimelineSteps, { stubs });

  test('Toggles the expanded state when expand button is clicked', async () => {
    const wrapper = factory({
      propsData: {
        steps,
      },
    });

    const firstAlarmTimelineStep = selectAlarmTimelineSteps(wrapper).at(0);
    firstAlarmTimelineStep.vm.$emit('expand', true);

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
    secondAlarmTimelineStep.vm.$emit('expand', true);

    await flushPromises();

    const firstAlarmTimelineStepExpand = selectAlarmTimelineStepExpands(wrapper).at(0);
    const secondAlarmTimelineStepExpand = selectAlarmTimelineStepExpands(wrapper).at(1);

    expect(firstAlarmTimelineStepExpand.props('expanded')).toBeFalsy();
    expect(secondAlarmTimelineStepExpand.props('expanded')).toBeTruthy();
  });
});
