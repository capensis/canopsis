<template lang="pug">
  c-action-btn(
    :icon="icon",
    :color="color",
    :tooltip="tooltip",
    :disabled="disabled || downloading",
    @click="downloadContent"
  )
</template>

<script>
import { saveJsonFile, saveCsvFile, saveTextFile, saveFile } from '@/helpers/file/files';

/**
 * @example
 *   c-download-btn(value="Text content", type="txt", name="txt_file") -> txt_file.txt
 *   c-download-btn(:value="{ json: 'json' }", type="json", name="json_file") -> json_file.json
 *   c-download-btn(value="csv", type="csv", name="csv_file") -> csv_file.csv
 *   c-download-btn(value="Text", name="without_type_file") -> without_type_file.txt
 *   c-download-btn(:value="blob", name="custom_type_file", type="log") -> custom_type_file.log
 */
export default {
  props: {
    value: {
      type: [String, Object, Number, Blob],
      required: true,
    },
    icon: {
      type: String,
      default: 'file_download',
    },
    type: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default() {
        return this.$t('common.download');
      },
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      downloading: false,
    };
  },
  methods: {
    download() {
      const handler = {
        json: saveJsonFile,
        csv: saveCsvFile,
        txt: saveTextFile,
      }[this.type];

      const name = this.name || new Date().toLocaleDateString();

      if (handler) {
        return handler(this.value, name);
      }

      const blob = this.value instanceof Blob
        ? this.value
        : new Blob([this.value], { type: 'charset=utf-8' });

      return saveFile(blob, [name, this.type].filter(Boolean).join('.'));
    },

    async downloadContent() {
      this.downloading = true;

      try {
        await this.download();
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || err.description || this.$t('errors.default') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>
