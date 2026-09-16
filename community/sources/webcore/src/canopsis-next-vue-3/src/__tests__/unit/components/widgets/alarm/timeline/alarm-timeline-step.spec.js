import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmTimelineStep from '@/components/widgets/alarm/timeline/alarm-timeline-step.vue';

const stubs = {
  'c-expand-btn': true,
  'alarm-timeline-step-icon': true,
  'alarm-timeline-step-title': true,
};

const selectExpandBtn = wrapper => wrapper.find('c-expand-btn-stub');

describe('alarm-timeline-step', () => {
  const factory = generateShallowRenderer(AlarmTimelineStep, { stubs });

  test('Expand button emits event', () => {
    const expanded = true;
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          steps: [{}],
          in_pbh: false,
          _t: 'some_step',
        },
        deep: false,
      },
    });

    const expandButton = selectExpandBtn(wrapper);

    expandButton.triggerCustomEvent('expand', expanded);

    expect(wrapper).toEmit('expand', expanded);
  });
});
