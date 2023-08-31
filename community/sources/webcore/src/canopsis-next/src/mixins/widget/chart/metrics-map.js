import { keyBy } from 'lodash';

export const widgetChartMetricsMap = {
  data() {
    return {
      widgetMetricsMap: {},
    };
  },
  created() {
    this.setWidgetMetricsMap();
  },
  methods: {
    setWidgetMetricsMap() {
      this.widgetMetricsMap = keyBy(this.widget.parameters?.metrics ?? [], 'metric');
    },
  },
};
