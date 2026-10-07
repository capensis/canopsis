import { generateShallowRenderer } from '@unit/utils/vue';

import { COLORS } from '@/config';
import { ALARM_LIST_STEPS } from '@/constants';

import { getAlarmStepColor, getAlarmStepIcon } from '@/helpers/entities/alarm/step/formatting';

import AlarmTimelineStep from '@/components/widgets/alarm/timeline/alarm-timeline-step.vue';

const stubs = {
  'c-expand-btn': true,
  'alarm-timeline-step-icon': true,
  'alarm-timeline-step-title': true,
};

const selectExpandBtn = wrapper => wrapper.find('c-expand-btn-stub');

describe('alarm-timeline-step', () => {
  const factory = generateShallowRenderer(AlarmTimelineStep, { stubs });

  test('Expand button emits event', async () => {
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

    await expandButton.triggerCustomEvent('update:modelValue', expanded);

    expect(wrapper).toEmit('expand', expanded);
  });

  test('Pause details are displayed separately from the comment', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          a: 'John Doe',
          m: 'Planned maintenance\nService: API',
          source_type: 'Component',
          source_name: 'API',
          pause_type: 'maintenance',
        },
      },
    });

    expect(wrapper.vm.pauseSource).toBe('Component API');
    expect(wrapper.vm.sanitizedStepMessage).toContain('Planned maintenance');
    expect(wrapper.vm.sanitizedStepMessage).toContain('Service: API');
    expect(wrapper.vm.pauseType).toBe('maintenance');
  });

  test('Inherited pause metadata comes from fields and leaves the comment untouched', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          in_pbh: true,
          m: 'Planned maintenance\nComponent: API',
          source_type: 'Component',
          source_name: 'API',
          pause_type: 'maintenance',
        },
      },
    });

    expect(wrapper.vm.pauseSource).toBe('Component API');
    expect(wrapper.vm.pauseType).toBe('maintenance');
    expect(wrapper.text()).toContain('Planned maintenance\nComponent: API');
  });

  test('Pause source metadata is read only from explicit step fields', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          in_pbh: false,
          m: 'Planned maintenance\nComponent: API',
          pause_type: 'maintenance',
        },
      },
    });

    expect(wrapper.vm.pauseSource).toBe('');
    expect(wrapper.vm.sanitizedStepMessage).toContain('Planned maintenance');
    expect(wrapper.vm.sanitizedStepMessage).toContain('Component: API');
  });

  test('Pause and resume steps use gray and blue colors', () => {
    expect(getAlarmStepColor(ALARM_LIST_STEPS.pause)).toBe(COLORS.state.pause);
    expect(getAlarmStepColor(ALARM_LIST_STEPS.resume)).toBe(COLORS.alarmSteps.activate);
  });

  test('Pause and resume steps use pause and play icons', () => {
    expect(getAlarmStepIcon(ALARM_LIST_STEPS.pause)).toBe('pause');
    expect(getAlarmStepIcon(ALARM_LIST_STEPS.resume)).toBe('play_arrow');
  });
});
