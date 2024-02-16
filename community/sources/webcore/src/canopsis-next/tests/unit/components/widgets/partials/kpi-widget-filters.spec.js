import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

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
  const snapshotFactory = generateRenderer(KpiWidgetFilters, { stubs });

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

    selectQuickDateIntervalField(wrapper).vm.$emit('input', newInterval);

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

    selectSamplingField(wrapper).vm.$emit('input', SAMPLINGS.month);

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

    selectFilterSelector(wrapper).vm.$emit('input', newFilters);

    expect(wrapper).toEmit('update:filters', newFilters);
  });

  test('Renders \'kpi-widget-filters\' with required props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widgetId: 'widget-id',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders \'kpi-widget-filters\' with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widgetId: 'id',
        interval: {},
        sampling: SAMPLINGS.month,
        filters: [{}, {}],
        userFilters: [{}, {}, {}],
        widgetFilters: [{}, {}, {}, {}],
        lockedFilter: 'locked-filter-id',
        minIntervalDate: 1633737600,
        showFilter: true,
        showSampling: true,
        showInterval: true,
        filterAddable: true,
        filterEditable: true,
        filterDisabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
