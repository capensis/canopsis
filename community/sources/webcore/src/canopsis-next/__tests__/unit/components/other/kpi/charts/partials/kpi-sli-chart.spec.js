import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import KpiSliChart from '@/components/other/kpi/charts/partials/kpi-sli-chart.vue';

const stubs = {
  'chart-export-actions': true,
};

const selectChartExportActions = wrapper => wrapper.find('chart-export-actions-stub');

describe('kpi-sli-chart', () => {
  const factory = generateRenderer(KpiSliChart, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Export csv event emitted', async () => {
    const wrapper = factory();

    await flushPromises();

    selectChartExportActions(wrapper).triggerCustomEvent('export:csv');

    expect(wrapper.emitted('export:csv')).toHaveLength(1);
  });

  it('Export png event emitted', async () => {
    const wrapper = factory();

    await flushPromises();

    selectChartExportActions(wrapper).triggerCustomEvent('export:png');

    expect(wrapper.emitted('export:png')).toHaveLength(1);
  });
});
