import ChartAnnotationPlugin from '@/externals/chart/plugins/annotation';

export const chartAnnotationMixin = {
  created() {
    this.addPlugin(ChartAnnotationPlugin);
  },
};

export default chartAnnotationMixin;
