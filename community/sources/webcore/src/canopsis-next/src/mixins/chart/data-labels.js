import { ChartDataLabels } from '@/externals/chart/plugins/data-labels';

export const chartDataLabelsMixin = {
  created() {
    this.addPlugin(ChartDataLabels);
  },
};
