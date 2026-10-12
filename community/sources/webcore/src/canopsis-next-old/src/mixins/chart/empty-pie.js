import { ChartEmptyPiePlugin } from '@/externals/chart/plugins/empty-pie';

export const chartEmptyPieMixin = {
  created() {
    this.addPlugin(ChartEmptyPiePlugin);
  },
};
