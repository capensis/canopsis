import { mockConsole, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { API_HOST, API_ROUTES, AVAILABILITY_FILENAME_PREFIX } from '@/config';
import { SAMPLINGS } from '@/constants';

import { saveFile } from '@/helpers/file/files';

import AvailabilityGraph from '@/components/other/availability/partials/availability-graph.vue';

vi.mock('@/helpers/file/files', () => ({
  saveFile: vi.fn(),
}));

const metrics = vi.hoisted(() => ({
  fetchAvailabilityHistoryWithoutStore: vi.fn().mockResolvedValue({ data: [] }),
  fetchAvailabilityHistoryExport: vi.fn().mockResolvedValue({
    _id: 'export-availability-history-id',
  }),
}));

vi.mock('@/stores/modules/metrics', async () => {
  const { defineStore } = await vi.importActual('pinia');

  return {
    useMetricsStore: defineStore('metrics', { actions: metrics }),
  };
});

const stubs = {
  'c-progress-overlay': true,
  'availability-history': true,
};

const selectAvailabilityHistory = wrapper => wrapper.find('availability-history-stub');

describe('availability-graph', () => {
  vi.useFakeTimers({ now: 1386435500000 });

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

  afterEach(() => {
    saveFile.mockClear();
    metrics.fetchAvailabilityHistoryWithoutStore.mockReset().mockResolvedValue({ data: [] });
    metrics.fetchAvailabilityHistoryExport.mockReset().mockResolvedValue({
      _id: 'export-availability-history-id',
    });
  });

  test('History fetched after mount', () => {
    factory({
      propsData: {
        availability,
        interval,
      },
    });

    expect(metrics.fetchAvailabilityHistoryWithoutStore).toHaveBeenCalledWith({
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
      propsData: {
        availability,
        interval: longInterval,
      },
    });

    metrics.fetchAvailabilityHistoryWithoutStore.mockClear();

    await selectAvailabilityHistory(wrapper).triggerCustomEvent('update:sampling', SAMPLINGS.day);

    expect(metrics.fetchAvailabilityHistoryWithoutStore).toHaveBeenCalledWith({
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
      propsData: {
        availability,
        interval,
      },
    });

    const data = new Blob();
    await selectAvailabilityHistory(wrapper).triggerCustomEvent('export:png', data);

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
      propsData: {
        availability,
        interval,
      },
    });

    const data = new Blob();
    await selectAvailabilityHistory(wrapper).triggerCustomEvent('export:png', data);

    expect(consoleMock.error).toHaveBeenCalledWith(error);
    expect($popups.error).toHaveBeenCalledWith({ text: errorMessage });
  });

  test('Widget exported after trigger export button', async () => {
    const windowOpenMock = vi.spyOn(window, 'open').mockReturnValue();

    const wrapper = factory({
      propsData: {
        availability,
        interval,
      },
    });

    await selectAvailabilityHistory(wrapper).triggerCustomEvent('export:csv');

    await flushPromises();

    expect(metrics.fetchAvailabilityHistoryExport).toHaveBeenCalledWith({
      params: defaultParams,
    });

    expect(windowOpenMock).toHaveBeenCalledWith(
      `${API_HOST}${API_ROUTES.metrics.exportMetric}/export-availability-history-id/download`,
      '_blank',
    );
    windowOpenMock.mockRestore();
  });

  test('Error popup showed after trigger export button with error', async () => {
    const windowOpenMock = vi.spyOn(window, 'open').mockReturnValue();
    metrics.fetchAvailabilityHistoryExport.mockRejectedValue(new Error('Export failed'));

    const wrapper = factory({
      propsData: {
        availability,
        interval,
      },
    });

    await selectAvailabilityHistory(wrapper).triggerCustomEvent('export:csv');

    await flushPromises();

    expect(metrics.fetchAvailabilityHistoryExport).toHaveBeenCalledWith({
      params: defaultParams,
    });

    expect($popups.error).toHaveBeenCalledWith({
      text: 'availability.popups.exportCSVFailed',
    });
    expect(windowOpenMock).not.toHaveBeenCalled();
    windowOpenMock.mockRestore();
  });
});
