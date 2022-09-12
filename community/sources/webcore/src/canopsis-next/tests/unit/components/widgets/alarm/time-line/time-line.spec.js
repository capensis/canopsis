import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { fakeAlarmDetails } from '@unit/data/alarm';

import TimeLine from '@/components/widgets/alarm/time-line/time-line.vue';

const localVue = createVueInstance();

const stubs = {
  'time-line-flag': true,
  'time-line-card': true,
  'c-pagination': {
    template: `
      <input class="c-pagination" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const snapshotStubs = {
  'time-line-flag': true,
  'time-line-card': true,
  'c-pagination': true,
};

const factory = (options = {}) => shallowMount(TimeLine, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(TimeLine, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('time-line', () => {
  const { steps } = fakeAlarmDetails();

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

  it('Renders `time-line` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `time-line` with isHtmlEnabled', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
        isHtmlEnabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
