import Faker from 'faker';
import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

import { KPI_SLI_GRAPH_DATA_TYPE, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiSliFilters from '@/components/other/kpi/charts/partials/kpi-sli-filters';

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
  const snapshotFactory = generateRenderer(KpiSliFilters, { stubs });

  it('Query changed after trigger a quick interval field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const quickIntervalField = wrapper.find('c-quick-date-interval-field-stub');

    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedInterval = {
      from: start,
      to: stop,
    };

    quickIntervalField.vm.$emit('input', {
      from: start,
      to: stop,
    });

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.interval).toEqual(expectedInterval);
    expect(eventData).toEqual({
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

    const samplingField = wrapper.find('c-sampling-field-stub');

    samplingField.vm.$emit('input', SAMPLINGS.month);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.sampling).toEqual(SAMPLINGS.month);
    expect(eventData).toEqual({
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

    const showModeField = wrapper.find('kpi-sli-show-mode-field-stub');

    showModeField.vm.$emit('input', KPI_SLI_GRAPH_DATA_TYPE.time);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.type).toEqual(KPI_SLI_GRAPH_DATA_TYPE.time);
    expect(eventData).toEqual({
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

    const filtersField = wrapper.find('c-filter-field-stub');

    filtersField.vm.$emit('input', expectedFilter);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.filter).toEqual(expectedFilter);
    expect(eventData).toEqual({
      ...initialQuery,
      filter: expectedFilter,
    });
  });

  it('Renders `kpi-sli-filters` with query', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        query: {
          sampling: SAMPLINGS.day,
          type: KPI_SLI_GRAPH_DATA_TYPE.percent,
          filter: null,
          interval: {
            from: QUICK_RANGES.last30Days.start,
            to: QUICK_RANGES.last30Days.stop,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
