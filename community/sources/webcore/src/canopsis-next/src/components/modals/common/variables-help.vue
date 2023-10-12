<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.variablesHelp.variables') }}
    template(#text="")
      v-layout(v-if="config.exportEntity", justify-end)
        v-btn(color="primary", small, @click="exportOriginal")
          v-icon(left) file_download
          span {{ $t('common.exportToJson') }}
      v-treeview(:items="config.variables", item-key="name")
        template(#prepend="{ item }")
          div.caption.font-italic(v-if="item.isArray") {{`(${$t('common.variableTypes.array')})`}}

        template(#label="{ item, leaf }")
          div {{ item.name }}
            span.pl-1(v-if="leaf") :
              c-ellipsis.pl-1.d-inline-block.text--secondary.body-1.pre-wrap(:text="String(item.value)")
            span.pl-1(v-else-if="!leaf && !(item.children && item.children.length)") :
              div.pl-1.d-inline-block.text--secondary.body-1.font-italic {{ $t('common.emptyObject') }}

        template(#append="{ leaf, item }")
          v-menu(bottom, left, offset-y)
            template(#activator="{ on }")
              v-tooltip(left)
                template(#activator="{ on: tooltipOn }")
                  v-btn(v-on="{ ...tooltipOn, ...on }", icon)
                    v-icon save_alt
                span {{ getTooltipContent(leaf, item) }}

            v-list(dense)
              v-list-tile(v-if="leaf", @click="copyPathToClipboard(item.path)")
                v-list-tile-avatar
                  v-icon content_copy
                v-list-tile-title {{ $t('common.copyPathToClipboard') }}
              v-list-tile(v-else, @click="exportAsJson(item)")
                v-list-tile-avatar
                  v-icon(size="24") $vuetify.icons.json
                v-list-tile-title {{ $t('common.exportToJson') }}
              v-list-tile(v-if="item.original", @click="exportAsPdf(item.original, config.exportPdfTemplate)")
                v-list-tile-avatar
                  v-icon(size="24") $vuetify.icons.pdf
                v-list-tile-title {{ $t('common.exportToPdf') }}
</template>

<script>
import { DATETIME_FORMATS, MODALS } from '@/constants';

import { saveJsonFile } from '@/helpers/file/files';
import { convertTreeviewToObject } from '@/helpers/treeview';
import { convertDateToString, getNowTimestamp } from '@/helpers/date/date';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { widgetActionsPanelAlarmExportPdfMixin } from '@/mixins/widget/actions-panel/alarm-export-pdf';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.variablesHelp,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    widgetActionsPanelAlarmExportPdfMixin,
  ],
  data() {
    return {
      exportAlarmToPdfPending: false,
    };
  },
  methods: {
    getTooltipContent(leaf, item) {
      if (item.original) {
        return leaf
          ? this.$t('common.copyFieldPathOrExportFieldToPdf', { field: item.name })
          : this.$t('common.exportFieldToPdfOrJson', { field: item.name });
      }

      return leaf
        ? this.$t('common.copyFieldPath', { field: item.name })
        : this.$t('common.exportFieldToJson', { field: item.name });
    },

    getNowDateString() {
      return convertDateToString(
        getNowTimestamp(),
        DATETIME_FORMATS.long,
      );
    },

    exportAsJson(item) {
      const object = convertTreeviewToObject(item);

      saveJsonFile(object, `${item.name}-${this.getNowDateString()}`);
    },

    exportOriginal() {
      saveJsonFile(this.config.exportEntity, `${this.config.exportEntityName}-${this.getNowDateString()}`);
    },

    async exportAsPdf(alarm, template) {
      this.exportAlarmToPdfPending = true;

      await this.exportAlarmToPdf(alarm, template);

      this.exportAlarmToPdfPending = false;
    },

    copyPathToClipboard(path) {
      this.$copyText(path)
        .then(this.onSuccessCopied)
        .catch(this.onErrorCopied);
    },

    onSuccessCopied() {
      this.$popups.success({ text: this.$t('success.pathCopied') });
    },

    onErrorCopied() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
