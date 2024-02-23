import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CAlarmActionChip from '@/components/common/chips/c-alarm-action-chip.vue';

const stubs = {
  'v-icon': {
    template: '<span v-on="$listeners" class="v-icon"></span>',
  },
};

const selectCloseIcon = wrapper => wrapper.find('.v-icon');

describe('c-alarm-action-chip', () => {
  const factory = generateShallowRenderer(CAlarmActionChip, { stubs });
  const snapshotFactory = generateRenderer(CAlarmActionChip);

  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = factory({
      propsData: {
        closable: true,
      },
    });

    const closeIcon = selectCloseIcon(wrapper);

    closeIcon.trigger('click');

    expect(wrapper).toHaveBeenEmit('close');
  });

  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
