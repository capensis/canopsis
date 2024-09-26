import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockPopups } from '@unit/utils/mock-hooks';

import EntityGantt from '@/components/widgets/alarm/entity-gantt/entity-gantt.vue';

const stubs = {
  'c-progress-overlay': true,
  'junit-gantt-chart': true,
};

const selectJunitGanttChart = wrapper => wrapper.find('junit-gantt-chart-stub');

describe('entity-gantt', () => {
  const $popups = mockPopups();
  const fetchItemGanttIntervalsWithoutStore = jest.fn().mockResolvedValue({
    meta: {},
    data: [],
  });
  const testSuiteModule = {
    name: 'testSuite/entityGantt',
    actions: {
      fetchItemGanttIntervalsWithoutStore,
    },
  };
  const store = createMockedStoreModules([
    testSuiteModule,
  ]);

  const ganttIntervals = [{}, {}];
  const ganttIntervalsMeta = {
    total_count: 2,
  };
  const alarm = {
    entity: {
      _id: 'entity-id',
    },
  };

  const factory = generateShallowRenderer(EntityGantt, { stubs });
  const snapshotFactory = generateRenderer(EntityGantt, { stubs });

  afterEach(() => {
    fetchItemGanttIntervalsWithoutStore.mockClear();
  });

  it('Gantt intervals fetched after mount', async () => {
    const wrapper = factory({
      store,
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

    expect(fetchItemGanttIntervalsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: alarm.entity._id,
        params: {
          limit: newQuery.itemsPerPage,
          page: newQuery.page,
        },
      },
    );
  });

  it('Gantt intervals fetched after change the query', async () => {
    fetchItemGanttIntervalsWithoutStore.mockResolvedValueOnce({
      meta: ganttIntervalsMeta,
      data: ganttIntervals,
    });
    factory({
      store,
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect(fetchItemGanttIntervalsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: alarm.entity._id,
        params: {
          limit: 10,
          page: 1,
        },
      },
    );
  });

  it('Error popup showed after fetch gantt intervals with error message', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

    const error = {
      message: 'Error message',
    };
    fetchItemGanttIntervalsWithoutStore.mockRejectedValueOnce(error);
    factory({
      store,
      propsData: {
        alarm,
      },
      mocks: {
        $popups,
      },
    });

    await flushPromises();

    expect($popups.error).toBeCalledWith({ text: error.message });
    consoleErrorSpy.mockClear();
  });

  it('Error popup showed after fetch gantt intervals with error description', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

    const error = {
      description: 'Error description',
    };
    fetchItemGanttIntervalsWithoutStore.mockRejectedValueOnce(error);
    factory({
      store,
      propsData: {
        alarm,
      },
      mocks: {
        $popups,
      },
    });

    await flushPromises();

    expect($popups.error).toBeCalledWith({ text: error.description });
    consoleErrorSpy.mockClear();
  });

  it('Error popup showed after fetch gantt intervals with error', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

    fetchItemGanttIntervalsWithoutStore.mockRejectedValueOnce({});
    factory({
      store,
      propsData: {
        alarm,
      },
      mocks: {
        $popups,
      },
    });

    await flushPromises();

    expect($popups.error).toBeCalledWith({ text: 'Something went wrong...' });
    consoleErrorSpy.mockClear();
  });

  it('Renders `entity-gantt` with required props', async () => {
    fetchItemGanttIntervalsWithoutStore
      .mockResolvedValueOnce({
        meta: ganttIntervalsMeta,
        data: ganttIntervals,
      });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm,
      },
    });

    expect(wrapper).toMatchSnapshot();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
