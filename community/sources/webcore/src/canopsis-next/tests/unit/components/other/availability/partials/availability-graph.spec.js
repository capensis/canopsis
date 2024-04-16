import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAvailabilityModule, createMockedStoreModules } from '@unit/utils/store';
import { mockConsole, mockPopups } from '@unit/utils/mock-hooks';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE, EXPORT_STATUSES, SAMPLINGS } from '@/constants';
import { API_HOST, API_ROUTES, AVAILABILITY_FILENAME_PREFIX } from '@/config';

import { saveFile } from '@/helpers/file/files';

import AvailabilityGraph from '@/components/other/availability/partials/availability-graph.vue';

jest.mock('@/helpers/file/files', () => ({
  saveFile: jest.fn(),
}));

const stubs = {
  'c-progress-overlay': true,
  'availability-history': true,
};

const selectAvailabilityHistory = wrapper => wrapper.find('availability-history-stub');

describe('availability-graph', () => {
  jest.useFakeTimers({ now: 1386435500000 });

  const $popups = mockPopups();
  const consoleMock = mockConsole();

  const availability = {
    uptime_duration: 3142,
    downtime_duration: 7101,
    uptime_share: '30.67',
    downtime_share: 69.33,
    entity: {
      _id: 'availability-entity-id-1',
    },
  };
  const interval = {
    from: 1384435500000,
    to: 1384435500000,
  };

  const {
    availabilityModule,
    exportAvailabilityData,
    fetchAvailabilityHistoryWithoutStore,
    fetchAvailabilityHistoryExport,
  } = createAvailabilityModule();
  const store = createMockedStoreModules([
    availabilityModule,
  ]);

  const defaultParams = {
    _id: availability.entity._id,
    sampling: SAMPLINGS.hour,
    from: interval.from,
    to: interval.to,
  };

  const factory = generateShallowRenderer(AvailabilityGraph, {
    stubs,
    mocks: {
      $popups,
    },
  });
  const snapshotFactory = generateRenderer(AvailabilityGraph, {
    stubs,
    mocks: {
      $popups,
    },
  });

  afterEach(() => {
    saveFile.mockClear();
  });

  test('History fetched after mount', () => {
    factory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    expect(fetchAvailabilityHistoryWithoutStore).toBeDispatchedWith({
      params: {
        _id: availability.entity._id,
        sampling: SAMPLINGS.hour,
        from: interval.from,
        to: interval.to,
      },
    });
  });

  test('History fetched after trigger update:sampling on availability history', async () => {
    const longInterval = {
      from: 1384425500000,
      to: 1384435500000,
    };

    const wrapper = factory({
      store,
      propsData: {
        availability,
        interval: longInterval,
      },
    });

    fetchAvailabilityHistoryWithoutStore.mockClear();

    await selectAvailabilityHistory(wrapper).triggerCustomEvent(
      'update:sampling',
      SAMPLINGS.day,
    );

    expect(fetchAvailabilityHistoryWithoutStore).toBeDispatchedWith({
      params: {
        _id: availability.entity._id,
        sampling: SAMPLINGS.day,
        from: longInterval.from,
        to: longInterval.to,
      },
    });
  });

  test('History exported as png after trigger export:png on availability history', async () => {
    const wrapper = factory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    const data = new Blob();
    await selectAvailabilityHistory(wrapper).triggerCustomEvent(
      'export:png',
      data,
    );

    expect(saveFile).toHaveBeenCalledWith(
      data,
      `${AVAILABILITY_FILENAME_PREFIX}-14/11/2013-14/11/2013-${SAMPLINGS.hour}`,
    );
  });

  test('Error popup showed after trigger export:png with error on availability history', async () => {
    const errorMessage = 'message';
    const error = new Error(errorMessage);
    saveFile.mockRejectedValue(error);

    const wrapper = factory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    const data = new Blob();
    await selectAvailabilityHistory(wrapper).triggerCustomEvent(
      'export:png',
      data,
    );

    expect(consoleMock.error).toHaveBeenCalledWith(error);
    expect($popups.error).toHaveBeenCalledWith({ text: errorMessage });
  });

  test('Widget exported after trigger export button', async () => {
    const windowOpenMock = jest.spyOn(window, 'open').mockReturnValue();

    const wrapper = factory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    await selectAvailabilityHistory(wrapper).triggerCustomEvent(
      'export:csv',
    );

    await flushPromises();

    expect(fetchAvailabilityHistoryExport).toBeDispatchedWith(
      {
        params: defaultParams,
      },
    );

    await flushPromises();

    jest.runAllTimers();

    expect(fetchAvailabilityHistoryExport).toBeDispatchedWith({
      params: defaultParams,
    });

    await flushPromises();

    expect(windowOpenMock).toHaveBeenCalledWith(
      `${API_HOST}${API_ROUTES.metrics.exportMetric}/${exportAvailabilityData._id}/download`,
      '_blank',
    );
  });

  test('Error popup showed after trigger export button with error', async () => {
    const windowOpenMock = jest.spyOn(window, 'open').mockReturnValue();
    const exportFailedAvailabilityData = {
      _id: 'export-availability-history-id',
      status: EXPORT_STATUSES.failed,
    };

    fetchAvailabilityHistoryExport.mockRejectedValue(exportFailedAvailabilityData);

    const wrapper = factory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    await selectAvailabilityHistory(wrapper).triggerCustomEvent(
      'export:csv',
    );

    await flushPromises();

    expect(fetchAvailabilityHistoryExport).toBeDispatchedWith({
      params: defaultParams,
    });

    await flushPromises();

    jest.runAllTimers();

    expect(fetchAvailabilityHistoryExport).toBeDispatchedWith({
      params: defaultParams,
    });

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Failed to export availabilities in CSV format',
    });
    expect(windowOpenMock).not.toHaveBeenCalled();
  });

  test('Renders `availability-graph` with default props.', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        availability,
        interval,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-graph` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        availability,
        interval,
        defaultShowType: AVAILABILITY_SHOW_TYPE.duration,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
