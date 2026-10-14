import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { KPI_SLI_GRAPH_DATA_TYPE, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiSliFilters from '@/components/other/kpi/charts/partials/kpi-sli-filters.vue';

const stubs = {
  'c-quick-date-interval-field': true,
  'c-filter-field': true,
  'c-sampling-field': true,
  'kpi-sli-show-mode-field': true,
};

describe('kpi-sli-filters', () => {
  const nowTimestamp = 1386435600000;
  const initialQuery = {
    sampling: SAMPLINGS.day,
    type: KPI_SLI_GRAPH_DATA_TYPE.percent,
    filter: null,
    interval: {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    },
  };

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiSliFilters, { stubs });

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

  it('Query changed after trigger a sampling field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateQueryField('sampling', SAMPLINGS.month);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      sampling: SAMPLINGS.month,
    });
  });

  it('Query changed after trigger a show mode field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.vm.updateQueryField('type', KPI_SLI_GRAPH_DATA_TYPE.time);

    await flushPromises();

    expect(wrapper).toEmit('update:query', {
      ...initialQuery,
      type: KPI_SLI_GRAPH_DATA_TYPE.time,
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
});
