import { ChartLimitedSegmentPlugin } from '@/externals/chart/plugins/limited-segment';

export const chartLimitedSegmentMixin = {
  created() {
    this.addPlugin(ChartLimitedSegmentPlugin);
  },
};
