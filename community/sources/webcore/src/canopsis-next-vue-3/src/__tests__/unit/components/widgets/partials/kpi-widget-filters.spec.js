import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';

const stubs = {
  'c-quick-date-interval-field': true,
  'c-sampling-field': true,
  'filter-selector': true,
  'filters-list-btn': true,
};

const selectQuickDateIntervalField = wrapper => wrapper.find('c-quick-date-interval-field-stub');
const selectSamplingField = wrapper => wrapper.find('c-sampling-field-stub');
const selectFilterSelector = wrapper => wrapper.find('filter-selector-stub');

describe('kpi-widget-filters', () => {
  const factory = generateShallowRenderer(KpiWidgetFilters, { stubs });

  test('Date interval changed after trigger quick date interval field', () => {
    const wrapper = factory({
      propsData: {
        widgetId: 'widgetId',
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
        showInterval: true,
      },
    });

    const newInterval = {
      from: Faker.datatype.number(),
      to: Faker.datatype.number(),
    };

    selectQuickDateIntervalField(wrapper).triggerCustomEvent('input', newInterval);

    expect(wrapper).toEmit('update:interval', newInterval);
  });

  test('Sampling changed after trigger sampling field', () => {
    const wrapper = factory({
      propsData: {
        widgetId: 'widgetId',
        sampling: SAMPLINGS.day,
        showSampling: true,
      },
    });

    selectSamplingField(wrapper).triggerCustomEvent('input', SAMPLINGS.month);

    expect(wrapper).toEmit('update:sampling', SAMPLINGS.month);
  });

  test('Filters changed after trigger filters selector', () => {
    const wrapper = factory({
      propsData: {
        widgetId: 'widgetId',
        showFilter: true,
      },
    });

    const newFilters = [Faker.datatype.string(), Faker.datatype.string()];

    selectFilterSelector(wrapper).triggerCustomEvent('input', newFilters);

    expect(wrapper).toEmit('update:filters', newFilters);
  });
});
