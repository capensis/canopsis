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

  test('Step message is sanitized and linkified', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          a: 'John Doe',
          m: 'See https://example.com for details',
        },
      },
    });

    expect(wrapper.vm.sanitizedStepMessage).toContain('<a href="https://example.com"');
    expect(wrapper.vm.sanitizedStepMessage).toContain('for details');
  });

  test('Plain step message is rendered unchanged when HTML is disabled', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          in_pbh: true,
          m: 'Planned maintenance\nWindow: API',
        },
      },
    });

    expect(wrapper.text()).toContain('Planned maintenance\nWindow: API');
  });

  test('Inherited source message is rendered exactly as received', () => {
    const message = 'Planned maintenance\nComponent: API';
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          m: message,
        },
      },
    });

    expect(wrapper.text()).toContain(message);
    expect(wrapper.text()).not.toContain('Inherited from:');
  });

  test('HTML step message is rendered through sanitized content', () => {
    const wrapper = factory({
      propsData: {
        step: {
          t: new Date(),
          _t: ALARM_LIST_STEPS.pause,
          in_pbh: false,
          m: '<strong>Planned maintenance</strong><script>alert(1)</script>',
        },
        isHtmlEnabled: true,
      },
    });

    expect(wrapper.vm.sanitizedStepMessage).toContain('<strong>Planned maintenance</strong>');
    expect(wrapper.vm.sanitizedStepMessage).not.toContain('<script>');
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
