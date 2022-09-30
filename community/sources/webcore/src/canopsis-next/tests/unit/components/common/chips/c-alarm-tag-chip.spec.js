import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CAlarmTagChip from '@/components/common/chips/c-alarm-tag-chip.vue';

const localVue = createVueInstance();

const stubs = {
  'v-icon': {
    template: '<span v-on="$listeners" class="v-icon"></span>',
  },
};

const factory = (options = {}) => shallowMount(CAlarmTagChip, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CAlarmTagChip, {
  localVue,

  ...options,
});

const selectCloseIcon = wrapper => wrapper.find('.v-icon');

describe('c-alarm-tag-chip', () => {
  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = factory({
      propsData: {
        close: true,
      },
    });

    const closeIcon = selectCloseIcon(wrapper);

    closeIcon.trigger('click');

    expect(wrapper).toEmit('close');
  });

  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-alarm-tag-chip` with custom props and slot', () => {
    const wrapper = snapshotFactory({
      propsData: {
        color: '#000',
        close: true,
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
