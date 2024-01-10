import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { fakeAlarmDetails } from '@unit/data/alarm';

import AlarmsTimeLine from '@/components/widgets/alarm/time-line/alarms-time-line.vue';

const stubs = {
  'alarms-time-line-flag': true,
  'alarms-time-line-card': true,
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

  it('Check pagination', () => {
    const page = 2;
    const wrapper = factory({
      propsData: {
        steps,
      },
    });

    const pagination = wrapper.find('.c-pagination');

    pagination.setValue(page);

    const updatePageEvents = wrapper.emitted('update:page');

    expect(updatePageEvents).toHaveLength(1);
    expect(updatePageEvents[0]).toEqual([page]);
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
