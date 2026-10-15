import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ALARM_METRIC_PARAMETERS, KPI_RATING_CRITERIA, QUICK_RANGES } from '@/constants';

import KpiRatingFilters from '@/components/other/kpi/charts/partials/kpi-rating-filters.vue';

const stubs = {
  'c-quick-date-interval-field': true,
  'c-filter-field': true,
  'kpi-rating-criteria-field': true,
  'kpi-rating-metric-field': true,
  'c-items-per-page-field': true,
};

describe('kpi-rating-filters', () => {
  const nowTimestamp = 1386435600000;
  const initialQuery = {
    filter: null,
    criteria: undefined,
    metric: ALARM_METRIC_PARAMETERS.ackAlarms,
    itemsPerPage: 5,
    interval: {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    },
  };

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiRatingFilters, { stubs });

  it('Query changed after trigger a quick interval field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedInterval = {
      from: start,
      to: stop,
    };

    wrapper.vm.updateQueryField('interval', expectedInterval);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      interval: expectedInterval,
    });
  });

  it('Query changed after trigger a filters field', async () => {
    const expectedFilter = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateQueryField('filter', expectedFilter);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      filter: expectedFilter,
    });
  });

  it('Query changed after trigger a criteria field', async () => {
    const ratingSetting = {
      id: 1,
      label: KPI_RATING_CRITERIA.role,
    };
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateCriteria(ratingSetting);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      criteria: ratingSetting,
    });
  });

  it('Query changed after change criteria to entity criteria', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateCriteria(KPI_RATING_CRITERIA.impactLevel);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      metric: ALARM_METRIC_PARAMETERS.createdAlarms,
      criteria: KPI_RATING_CRITERIA.impactLevel,
    });
  });

  it('Query changed after trigger a parameters field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateQueryField('metric', ALARM_METRIC_PARAMETERS.cancelAckAlarms);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      metric: ALARM_METRIC_PARAMETERS.cancelAckAlarms,
    });
  });
});
