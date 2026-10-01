import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmStateCounters from '@/components/widgets/service-weather/alarm-state-counters.vue';

const CTooltipStub = {
  name: 'c-tooltip',
  props: ['theme'],
  template: '<div><slot /></div>',
};

describe('alarm-state-counters', () => {
  const factory = generateShallowRenderer(AlarmStateCounters, {
    stubs: {
      'c-tooltip': CTooltipStub,
    },
  });

  test.each([
    { dark: false, theme: 'light' },
    { dark: true, theme: 'dark' },
  ])('uses the $theme theme for its tooltip', ({ dark, theme }) => {
    const wrapper = factory({
      parentComponent: {
        provide: {
          $system: { dark },
        },
      },
    });

    expect(wrapper.findComponent(CTooltipStub).props('theme')).toBe(theme);
  });
});
