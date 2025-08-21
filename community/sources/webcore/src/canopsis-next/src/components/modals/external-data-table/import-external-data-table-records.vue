<template>
  <v-form class="import-external-data-table-records-form" @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-data-table-general-info-form :form="config.externalDataTable" />
        <c-csv-separator-field v-model="separator" :disabled="!!activeImportFileId" />
        <file-drag-selector
          :file-type-label="$t('common.fileSelector.fileTypes.csv')"
          :max-file-size="fileImportMaxSizeInKb"
          :loading="uploading"
          accept=".csv"
          required
          @change="chooseFile"
        >
          <template #selected="{ files, clear }">
            <v-layout class="position-relative" column>
              <c-progress-overlay :pending="uploading" />
              <v-layout
                v-for="file in files"
                :key="file.name"
                class="filename"
                justify-space-between
                align-center
              >
                <strong>
                  <v-icon>upload_file</v-icon>
                  {{ file.name }}
                </strong>
                <c-action-btn
                  type="delete"
                  @click="clear"
                />
              </v-layout>
            </v-layout>
          </template>
          <template #label="{ on }">
            <c-progress-overlay :pending="uploading" />
            <span>
              <p class="text-subtitle-2">
                {{ $t('common.fileSelector.dragAndDrop.label') }}
                <a v-on="on">
                  {{ $t('common.fileSelector.dragAndDrop.labelAction') }}
                </a>
                {{ $t('common.fileSelector.fileTypes.csv') }}
                ({{ $t('common.fileSelector.fileSizeMb', { size: fileImportMaxSizeInMb }) }})
              </p>
              <p class="text-subtitle-2">{{ $t('externalData.importFileDescription') }}</p>
            </span>
            <v-btn
              v-if="hasStructure"
              :loading="downloading"
              class="ml-3"
              color="primary"
              outlined
              @click="exportTableStructure"
            >
              {{ $t('externalData.exportTableStructure') }}
            </v-btn>
          </template>
        </file-drag-selector>
        <template v-if="activeImportFileId">
          <span class="text-subtitle-2">{{ $t('common.preview') }}</span>
          <external-data-table-records-list
            v-model="form"
            :records="records"
            :pending="pending"
            :options="options"
            :total-items="meta.total_count"
            :separator="separator"
            has-structure
            @update:options="updateOptions"
          />
        </template>
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.import') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import {
  CSV_SEPARATORS,
  EXTERNAL_METRIC_UNITS,
  IMPORT_STATUSES,
  MODALS,
  VALIDATION_DELAY,
} from '@/constants';

import { convertFileSizeToUnit } from '@/helpers/file/size';
import { saveCsvFile } from '@/helpers/file/files';
import { externalDataTableColumnsConfigToForm } from '@/helpers/entities/external-data-table/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useValidator } from '@/hooks/validator/validator';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useFilePolling } from '@/hooks/polling';
import { useInfo } from '@/hooks/store/modules/info';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';
import { useExternalDataTableImport } from '@/hooks/store/modules/external-data-table-import';

import ExternalDataTableRecordsList
  from '@/components/other/external-data-table/partials/external-data-table-records-list.vue';
import ExternalDataTableGeneralInfoForm
  from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';
import FileDragSelector from '@/components/forms/fields/file-drag-selector.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importExternalDataTableRecords,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ExternalDataTableRecordsList,
    FileDragSelector,
    ExternalDataTableGeneralInfoForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const validator = useValidator();

    const { fileImportMaxSize } = useInfo();
    const { fetchExternalDataTableSchema } = useExternalDataTable();
    const {
      createExternalDataTableImport,
      fetchExternalDataTableImportStatus,
      fetchExternalDataTableImportData,
      completeExternalDataTableImport,
    } = useExternalDataTableImport();

    const separator = ref(CSV_SEPARATORS.comma);
    const activeImportFileId = ref(null);
    const form = ref({});

    const fileImportMaxSizeInKb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.kilobyte)
    ));

    const fileImportMaxSizeInMb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.megabyte)
    ));

    const hasStructure = computed(() => !!config.value.externalDataTable?.columns?.length);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form: separator,
      method: async () => {
        await completeExternalDataTableImport({
          id: activeImportFileId.value,
          data: {
            column_types: Object.values(form.value),
          },
        });
        await config.value.afterSubmit?.(activeImportFileId.value);

        close();
      },
    });

    const title = computed(() => config.value.title || t('modals.createExternalDataTableRecord.create.title'));

    const {
      data: records,
      meta,
      pending,
      resetQuery,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: rest => fetchExternalDataTableImportData({
        ...rest,
        id: activeImportFileId.value,
      }),
    });

    /**
     * Handles file generation and download for technical metrics export.
     */
    const { poll: importFile } = useFilePolling({
      createHandler: createExternalDataTableImport,
      fetchHandler: fetchExternalDataTableImportStatus,
      completedStatus: IMPORT_STATUSES.completed,
      failedStatus: IMPORT_STATUSES.failed,
    });

    const {
      pending: uploading,
      handler: chooseFile,
    } = usePendingHandler(async ([file] = []) => {
      try {
        resetQuery();

        activeImportFileId.value = null;

        if (!file) {
          return;
        }

        validator.errors.clear();

        const data = {
          separator: separator.value,
          file,
        };

        const {
          _id: id,
          column_configs: columnConfigs = [],
        } = await importFile({ id: config.value.externalDataTable._id, data });

        activeImportFileId.value = id;

        form.value = externalDataTableColumnsConfigToForm(columnConfigs, true);

        fetchList();
      } catch (err) {
        if (!err.file) {
          throw err;
        }

        validator.errors.add({
          field: 'file',
          msg: err.file,
        });
      }
    });

    const {
      pending: downloading,
      handler: exportTableStructure,
    } = usePendingHandler(async () => {
      const blob = await fetchExternalDataTableSchema({ id: config.value.externalDataTable._id });

      return saveCsvFile(blob, `external-data-table-${config.value.externalDataTable._id}-structure`);
    });

    return {
      config,
      fileImportMaxSize,
      fileImportMaxSizeInKb,
      fileImportMaxSizeInMb,
      separator,
      activeImportFileId,

      hasStructure,

      form,

      isDisabled,
      submitting,

      pending,
      meta,
      records,
      options,
      updateOptions,

      title,

      uploading,
      downloading,

      submit,
      close,
      chooseFile,
      exportTableStructure,
    };
  },
};
</script>

<style lang="scss">
.import-external-data-table-records-form {
  .drag-zone {
    padding: 40px 12px !important;
    justify-content: center;

    p {
      margin: 0;
    }
  }

  .filename {
    padding: 12px;
    border: 2px solid;
    border-radius: 15px;

    .theme--light & {
      border-color: var(--v-application-background-darken2);

      .v-icon {
        color: var(--v-application-background-darken2);
      }
    }

    .theme--dark & {
      border-color: var(--v-application-background-lighten3);

      .v-icon {
        color: var(--v-application-background-lighten4);
      }
    }
  }
}
</style>
