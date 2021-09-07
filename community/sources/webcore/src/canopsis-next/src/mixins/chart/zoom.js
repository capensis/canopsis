import { ChartPluginZoom } from '@/externals/chart/plugins/zoom';

export const chartZoomMixin = {
  created() {
    this.addPlugin(ChartPluginZoom);
  },
};
