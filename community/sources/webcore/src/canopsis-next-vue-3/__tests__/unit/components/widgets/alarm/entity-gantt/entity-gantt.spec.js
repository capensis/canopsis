import { mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import EntityGantt from '@/components/widgets/alarm/entity-gantt/entity-gantt.vue';

const testSuite = vi.hoisted(() => ({
  fetchTestSuiteEntityGanttIntervalsWithoutStore: vi.fn().mockResolvedValue({
    meta: {},
    data: [],
  }),
}));

vi.mock('@/hooks/store/modules/test-suite', () => ({
  useTestSuite: () => testSuite,
}));

const stubs = {
  'c-progress-overlay': true,
  'junit-gantt-chart': true,
};

const selectJunitGanttChart = wrapper => wrapper.find('junit-gantt-chart-stub');

describe('entity-gantt', () => {
  const $popups = mockPopups();
  const ganttIntervals = [{}, {}];
  const ganttIntervalsMeta = {
    total_count: 2,
  };
  const alarm = {
    entity: {
      _id: 'entity-id',
    },
  };

  const factory = generateShallowRenderer(EntityGantt, { stubs, mocks: { $popups } });

  afterEach(() => {
    testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore
      .mockReset()
      .mockResolvedValue({ meta: {}, data: [] });
  });

  it('Gantt intervals fetched after mount', async () => {
    const wrapper = factory({
      propsData: {
        alarm,
      },
    });

    const junitGanttChart = selectJunitGanttChart(wrapper);

    const newQuery = {
      itemsPerPage: 15,
      page: 2,
    };

    junitGanttChart.triggerCustomEvent('update:query', newQuery);

    await flushPromises();

    expect(testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore).toHaveBeenCalledWith({
      id: alarm.entity._id,
      params: {
        limit: newQuery.itemsPerPage,
        page: newQuery.page,
      },
    });
  });

  it('Gantt intervals fetched after change the query', async () => {
    testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore.mockResolvedValueOnce({
      meta: ganttIntervalsMeta,
      data: ganttIntervals,
    });
    factory({
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect(testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore).toHaveBeenCalledWith({
      id: alarm.entity._id,
      params: {
        limit: 10,
        page: 1,
      },
    });
  });

  it('Error popup showed after fetch gantt intervals with error message', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    const error = {
      message: 'Error message',
    };
    testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore.mockRejectedValueOnce(error);
    factory({
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({ text: error.message });
    consoleErrorSpy.mockRestore();
  });

  it('Error popup showed after fetch gantt intervals with error description', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    const error = {
      description: 'Error description',
    };
    testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore.mockRejectedValueOnce(error);
    factory({
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({ text: error.description });
    consoleErrorSpy.mockRestore();
  });

  it('Error popup showed after fetch gantt intervals with error', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    testSuite.fetchTestSuiteEntityGanttIntervalsWithoutStore.mockRejectedValueOnce({});
    factory({
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({ text: 'errors.default' });
    consoleErrorSpy.mockRestore();
  });
});
