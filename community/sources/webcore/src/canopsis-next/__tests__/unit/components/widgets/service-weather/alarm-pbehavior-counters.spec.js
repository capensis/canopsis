import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmPbehaviorCounters from '@/components/widgets/service-weather/alarm-pbehavior-counters.vue';

const CTooltipStub = {
  name: 'c-tooltip',
  props: ['theme'],
  template: '<div><slot /></div>',
};

describe('alarm-pbehavior-counters', () => {
  const factory = generateShallowRenderer(AlarmPbehaviorCounters, {
    stubs: {
      'c-tooltip': CTooltipStub,
    },
  });

  test.each([
    { dark: false, theme: 'light' },
    { dark: true, theme: 'dark' },
  ])('uses the $theme theme for its tooltips', ({ dark, theme }) => {
    const wrapper = factory({
      parentComponent: {
        provide: {
          $system: { dark },
        },
      },
      propsData: {
        counters: [
          {
            count: 1,
            type: {
              _id: 'selected',
              name: 'Selected',
              icon_name: 'check',
            },
          },
          {
            count: 2,
            type: {
              _id: 'other',
              name: 'Other',
            },
          },
        ],
        types: ['selected'],
      },
    });

    expect(wrapper.findAllComponents(CTooltipStub)).toHaveLength(2);
    wrapper.findAllComponents(CTooltipStub).forEach((tooltip) => {
      expect(tooltip.props('theme')).toBe(theme);
    });
  });
});
