import ChartAnnotationPlugin from '@/externals/chart/plugins/annotation';

export default {
  created() {
    this.addPlugin(ChartAnnotationPlugin);
  },
};
