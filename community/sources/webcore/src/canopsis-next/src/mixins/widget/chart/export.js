import { API_HOST, API_ROUTES, KPI_ALARM_METRICS_FILENAME_PREFIX } from '@/config';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';
import { exportMixinCreator } from '@/mixins/widget/export';

export const widgetChartExportMixinCreator = ({ createExport, fetchExport }) => ({
  mixins: [
    exportMixinCreator({ createExport, fetchExport }),
  ],
  data() {
    return {
      downloading: false,
    };
  },
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

    async exportMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    async exportMetricsAsCsv() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          data: this.getQuery(),
        });

        this.downloadFile(`${API_HOST}${API_ROUTES.metrics.exportMetric}/${fileData._id}/download`);
      } catch (err) {
        this.$popups.error({ text: this.$t('kpi.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
});
