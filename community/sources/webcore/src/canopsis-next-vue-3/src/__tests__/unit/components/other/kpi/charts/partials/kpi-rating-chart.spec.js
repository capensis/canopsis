import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import KpiRatingChart from '@/components/other/kpi/charts/partials/kpi-rating-chart.vue';

const stubs = {
  'chart-export-actions': true,
};

describe('kpi-rating-chart', () => {
  const factory = generateRenderer(KpiRatingChart, {
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

    wrapper.vm.$emit('export:csv');

    expect(wrapper.emitted('export:csv')).toHaveLength(1);
  });

  it('Export png event emitted', async () => {
    const wrapper = factory();

    await flushPromises();

    wrapper.vm.$emit('export:png');

    expect(wrapper.emitted('export:png')).toHaveLength(1);
  });
});
