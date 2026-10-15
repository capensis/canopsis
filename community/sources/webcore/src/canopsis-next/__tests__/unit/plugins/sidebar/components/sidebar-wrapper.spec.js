import { createButtonStub } from '@unit/stubs/button';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { SIDE_BARS } from '@/constants';

import SidebarWrapper from '@/components/layout/sidebar/sidebar-wrapper.vue';

const stubs = {
  'v-btn': createButtonStub('v-btn'),
  'v-navigation-drawer': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: `
      <div class="v-navigation-drawer">
        <input
          :checked="modelValue"
          type="checkbox"
          @change="$emit('update:modelValue', $event.target.checked)"
        />
        <slot />
      </div>
    `,
  },
};

describe('sidebar-wrapper', () => {
  const sidebar = {
    id: 'test-sidebar-id',
    name: SIDE_BARS.alarmSettings,
    config: {},
    hidden: false,
  };

  const factory = generateShallowRenderer(SidebarWrapper, {
    parentComponent: {
      provide: {
        $clickOutside: {
          call: vi.fn(),
        },
      },
    },

    stubs,
  });

  beforeEach(() => {
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((fn) => {
      fn(0);

      return 0;
    });
  });

  afterEach(() => {
    window.requestAnimationFrame.mockRestore();
  });

  it('Close emitted after drawer is hidden', async () => {
    const wrapper = factory({
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    const navigationDrawerInput = wrapper.find('.v-navigation-drawer input');

    await navigationDrawerInput.setChecked(false);

    expect(wrapper.emitted('close')).toHaveLength(1);
  });

  it('Close emitted after click on the overlay', async () => {
    const wrapper = factory({
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    await wrapper.find('.sidebar-wrapper-overlay').trigger('click');

    expect(wrapper.emitted('close')).toHaveLength(1);
  });
});
