import { KPI_ALARM_METRICS_FILENAME_PREFIX } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { metricsExportMixinCreator } from '@/mixins/widget/metrics/export';

export const widgetChartExportMixinCreator = ({ createExport, fetchExport }) => ({
  mixins: [
    metricsExportMixinCreator({ createExport, fetchExport }),
  ],
  methods: {
    getFileName() {
      const { from, to } = this.getIntervalQuery();

      const fromTime = convertDateToString(from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(to, DATETIME_FORMATS.short);

      return [
        KPI_ALARM_METRICS_FILENAME_PREFIX,
        this.widget.parameters.chart_title,
        fromTime,
        toTime,
        this.query.sampling,
      ]
        .filter(Boolean)
        .join('-');
    },
  },
});
