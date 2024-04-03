import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import AlarmsTimeLine from '@/components/widgets/alarm/time-line/alarms-time-line.vue';
import AlarmsTimeLineSteps from '@/components/widgets/alarm/time-line/alarms-time-line-steps.vue';

const stubs = {
  'alarms-time-line-flag': true,
  'alarms-time-line-card': true,
  'alarms-time-line-steps': AlarmsTimeLineSteps,
  'c-pagination': {
    template: `
      <input class="c-pagination" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const snapshotStubs = {
  'alarms-time-line-flag': true,
  'alarms-time-line-card': true,
  'c-pagination': true,
};

describe('alarms-time-line', () => {
  const { steps } = fakeAlarmDetails();

  const factory = generateShallowRenderer(AlarmsTimeLine, { stubs });
  const snapshotFactory = generateRenderer(AlarmsTimeLine, {
    stubs: snapshotStubs,
  });

  it('Check pagination', async () => {
    const page = 2;
    const wrapper = factory({
      propsData: {
        steps,
      },
    });

    await wrapper.find('.c-pagination').setValue(page);

    expect(wrapper).toEmit('update:page', page);
  });

  it('Renders `alarms-time-line` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line` with isHtmlEnabled', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
        isHtmlEnabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
