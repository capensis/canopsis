<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.variablesHelp.variables') }}</span>
    </template>
    <template #text="">
      <v-layout
        v-if="config.exportEntity"
        justify-end
      >
        <v-btn
          color="primary"
          small
          @click="exportOriginal"
        >
          <v-icon left>
            file_download
          </v-icon>
          <span>{{ $t('common.exportToJson') }}</span>
        </v-btn>
      </v-layout>
      <v-treeview
        :items="config.variables"
        item-key="name"
      >
        <template #prepend="{ item }">
          <div
            class="text-caption font-italic"
            v-if="item.isArray"
          >
            {{ `(${$t('common.variableTypes.array')})` }}
          </div>
        </template>
        <template #label="{ item, leaf }">
          <div>
            {{ item.name }}
            <span
              class="pl-1"
              v-if="leaf"
            >:
              <c-ellipsis
                class="pl-1 d-inline-block text--secondary text-body-1 pre-wrap"
                :text="String(item.value)"
              />
            </span>
            <span
              class="pl-1"
              v-else-if="!leaf && !(item.children && item.children.length)"
            >:
              <div class="pl-1 d-inline-block text--secondary text-body-1 font-italic">
                {{ $t('common.emptyObject') }}
              </div>
            </span>
          </div>
        </template>
        <template #append="{ leaf, item }">
          <v-menu
            bottom
            left
            offset-y
          >
            <template #activator="{ on }">
              <v-tooltip left>
                <template #activator="{ on: tooltipOn }">
                  <v-btn
                    v-on="{ ...tooltipOn, ...on }"
                    icon
                  >
                    <v-icon>{{ leaf ? 'content_copy' : 'save_alt' }}</v-icon>
                  </v-btn>
                </template>
                <span>{{ getTooltipContent(leaf, item) }}</span>
              </v-tooltip>
            </template>
            <v-list dense>
              <v-list-item
                v-if="leaf"
                @click="copyPathToClipboard(item.path)"
              >
                <v-list-item-avatar>
                  <v-icon>content_copy</v-icon>
                </v-list-item-avatar>
                <v-list-item-title>{{ $t('common.copyPathToClipboard') }}</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-else
                @click="exportAsJson(item)"
              >
                <v-list-item-avatar>
                  <v-icon size="24">
                    $vuetify.icons.json
                  </v-icon>
                </v-list-item-avatar>
                <v-list-item-title>{{ $t('common.exportToJson') }}</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-if="item.original"
                @click="exportAsPdf(item.original, config.exportPdfTemplate)"
              >
                <v-list-item-avatar>
                  <v-icon size="24">
                    $vuetify.icons.pdf
                  </v-icon>
                </v-list-item-avatar>
                <v-list-item-title>{{ $t('common.exportToPdf') }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </v-treeview>
    </template>
  </modal-wrapper>
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
