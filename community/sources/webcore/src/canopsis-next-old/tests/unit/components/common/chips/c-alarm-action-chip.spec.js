import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CChip from '@/components/common/chips/c-chip.vue';

const stubs = {
  'v-icon': {
    template: '<span v-on="$listeners" class="v-icon"></span>',
  },
};

const selectCloseIcon = wrapper => wrapper.find('.v-icon');

describe('c-chip', () => {
  const factory = generateShallowRenderer(CChip, { stubs });
  const snapshotFactory = generateRenderer(CChip);

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
