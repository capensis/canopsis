import { ChartBackgroundPlugin } from '@/externals/chart/plugins/background';

export const chartBackgroundMixin = {
  created() {
    this.addPlugin(ChartBackgroundPlugin);
  },
};
