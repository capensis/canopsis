import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import AlarmTimelineStep from '@/components/widgets/alarm/timeline/alarm-timeline-step.vue';

const stubs = {
  'c-expand-btn': true,
  'alarm-timeline-step-icon': true,
  'alarm-timeline-step-title': true,
};

const selectExpandBtn = wrapper => wrapper.find('c-expand-btn-stub');

describe('alarm-timeline-step', () => {
  const timestamp = 1386435600;
  const factory = generateShallowRenderer(AlarmTimelineStep, { stubs });
  const snapshotFactory = generateRenderer(AlarmTimelineStep, { stubs });

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

  test('Renders `alarm-timeline-step` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-timeline-step` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          t: timestamp,
          steps: [],
          in_pbh: false,
          _t: 'some_step',
        },
        deep: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-timeline-step` with deep and pbh correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          t: timestamp,
          steps: [],
          in_pbh: false,
          _t: 'some_step',
        },
        deep: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-timeline-step` with children correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          t: timestamp,
          steps: [{}],
          in_pbh: false,
          _t: 'some_step',
        },
        deep: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
