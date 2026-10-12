import { generateShallowRenderer } from '@unit/utils/vue';

import CAlarmExtraDetailsChip from '@/components/common/chips/c-alarm-extra-details-chip.vue';

describe('c-alarm-extra-details-chip', () => {
  const system = {
    theme: {
      colors: {
        icons: {
          ack: '#FFFFFF',
        },
      },
    },
  };
  const factory = generateShallowRenderer(CAlarmExtraDetailsChip, {
    global: { provide: { $system: system } },
  });

  it('uses the active theme color and calculates a contrasting icon color', () => {
    const wrapper = factory({ propsData: { icon: 'done', themeColor: 'ack' } });

    expect(wrapper.attributes('style')).toContain('background-color: rgb(255, 255, 255)');
    expect(wrapper.vm.preparedIconColor).toBe('black');
  });

  it('uses a light icon for a dark theme color', () => {
    const wrapper = factory({
      propsData: { icon: 'done', themeColor: 'ack' },
      global: {
        provide: {
          $system: { theme: { colors: { icons: { ack: '#000000' } } } },
        },
      },
    });

    expect(wrapper.attributes('style')).toContain('background-color: rgb(0, 0, 0)');
    expect(wrapper.vm.preparedIconColor).toBe('white');
  });
});
