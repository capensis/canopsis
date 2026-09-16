import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import AvailabilityLineChart from '@/components/other/availability/partials/availability-line-chart.vue';

const stubs = {
  'chart-export-actions': true,
};

const selectChartExportActions = wrapper => wrapper.find('chart-export-actions-stub');

describe('availability-line-chart', () => {
  const factory = generateRenderer(AvailabilityLineChart, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Export csv event emitted', async () => {
    const exportCsv = vi.fn();

    const wrapper = factory({
      listeners: {
        'export:csv': exportCsv,
      },
    });

    await flushPromises();

    selectChartExportActions(wrapper).triggerCustomEvent('export:csv');

    expect(exportCsv).toHaveBeenCalledTimes(1);
  });

  it('Export png event emitted', async () => {
    const exportPng = vi.fn();

    const wrapper = factory({
      listeners: {
        'export:png': exportPng,
      },
    });

    await flushPromises();

    selectChartExportActions(wrapper).triggerCustomEvent('export:png');

    expect(exportPng).toHaveBeenCalledTimes(1);
  });
});
