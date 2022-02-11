import { mount, createVueInstance } from '@unit/utils/vue';

import AlarmHeaderPriority from '@/components/widgets/alarm/headers-formatting/alarm-header-priority.vue';

const localVue = createVueInstance();

const stubs = {
  'c-help-icon': true,
};

const snapshotFactory = (options = {}) => mount(AlarmHeaderPriority, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-header-priority', () => {
  it('Renders `alarm-header-priority` without slot', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-header-priority` with slot', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: 'Default text slot',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
