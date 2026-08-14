<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ t('modals.variablesHelp.variables') }}</span>
    </template>

    <template #text>
      <c-layout v-if="config.exportEntity" justify="end" class="mb-2">
        <v-btn color="primary" size="small" @click="exportOriginal">
          <template #prepend>
            <v-icon icon="file_download" />
          </template>

          <span>{{ t('common.exportToJson') }}</span>
        </v-btn>
      </c-layout>

      <c-treeview :items="config.variables" item-key="path" item-title="name">
        <template #prepend="{ item }">
          <div v-if="item.isArray" class="text-caption font-italic text-medium-emphasis">
            {{ `(${t('common.variableTypes.array')})` }}
          </div>
        </template>

        <template #label="{ item, leaf }">
          <div>
            {{ item.name }}
            <span
              v-if="leaf"
              class="pl-1"
            >:
              <c-ellipsis
                :text="String(item.value)"
                class="pl-1 d-inline-block text-medium-emphasis text-body-1 pre-wrap variables-help__value"
              />
            </span>
            <span
              v-else-if="!leaf && !(item.children && item.children.length)"
              class="pl-1"
            >:
              <div class="pl-1 d-inline-block text-medium-emphasis text-body-1 font-italic">
                {{ t('common.emptyObject') }}
              </div>
            </span>
          </div>
        </template>

        <template #append="{ leaf, item }">
          <v-menu location="bottom end">
            <template #activator="{ props: menuProps }">
              <v-btn
                v-bind="menuProps"
                :icon="leaf ? 'content_copy' : 'save_alt'"
                size="small"
                variant="text"
                @click.stop
              />
            </template>

            <v-list density="compact">
              <v-list-item v-if="leaf" @click="copyPathToClipboard(item.path)">
                <template #prepend>
                  <v-icon icon="content_copy" />
                </template>

                <v-list-item-title>{{ t('common.copyPathToClipboard') }}</v-list-item-title>
              </v-list-item>

              <v-list-item v-else @click="exportAsJson(item)">
                <template #prepend>
                  <v-icon icon="$json" />
                </template>

                <v-list-item-title>{{ t('common.exportToJson') }}</v-list-item-title>
              </v-list-item>

              <v-list-item
                v-if="item.original"
                @click="exportAlarmToPdf(item.original, config.exportPdfTemplate)"
              >
                <template #prepend>
                  <v-icon icon="$pdf" />
                </template>

                <v-list-item-title>{{ t('common.exportToPdf') }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </c-treeview>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { DATETIME_FORMATS, MODALS } from '@/constants';

import { convertDateToString, getNowTimestamp } from '@/helpers/date/date';
import { saveJsonFile } from '@/helpers/file/files';
import { convertTreeviewToObject } from '@/helpers/treeview';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useAlarmExportPdf } from '@/hooks/widget/alarm-export-pdf';

export default {
  name: MODALS.variablesHelp,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config } = useInnerModal(props);
    const popups = usePopups();
    const { exportAlarmToPdf } = useAlarmExportPdf();

    const getNowDateString = () => convertDateToString(getNowTimestamp(), DATETIME_FORMATS.long);

    const exportAsJson = (item) => {
      const object = convertTreeviewToObject(item);

      saveJsonFile(object, `${item.name}-${getNowDateString()}`);
    };

    const exportOriginal = () => {
      saveJsonFile(
        config.value.exportEntity,
        `${config.value.exportEntityName}-${getNowDateString()}`,
      );
    };

    const copyPathToClipboard = async (path) => {
      try {
        await navigator.clipboard.writeText(path);

        popups.success({ text: t('success.pathCopied') });
      } catch {
        popups.error({ text: t('errors.default') });
      }
    };

    return {
      t,
      config,
      exportAlarmToPdf,
      exportAsJson,
      exportOriginal,
      copyPathToClipboard,
    };
  },
};
</script>

<style lang="scss" scoped>
.variables-help__value {
  max-width: 520px;
  vertical-align: bottom;
}
</style>
