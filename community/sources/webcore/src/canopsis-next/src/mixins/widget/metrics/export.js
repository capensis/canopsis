import { KPI_RATING_METRICS_FILENAME_PREFIX } from '@/config';

import { saveFile } from '@/helpers/file/files';
import { getExportMetricDownloadFileUrl } from '@/helpers/entities/metric/url';

import { exportMixinCreator } from '@/mixins/widget/export';

export const metricsExportMixinCreator = ({ createExport, fetchExport }) => ({
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
      return KPI_RATING_METRICS_FILENAME_PREFIX;
    },

    async exportMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    async exportMetricsAsCsv() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          data: this.getQuery(),
        });

        this.downloadFile(getExportMetricDownloadFileUrl(fileData._id));
      } catch (err) {
        this.$popups.error({ text: this.$t('kpi.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
});
