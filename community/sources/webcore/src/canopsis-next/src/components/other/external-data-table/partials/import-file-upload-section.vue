<template>
  <div>
    <c-csv-separator-field :value="separator" :disabled="!!activeImportFileId" @input="updateSeparator" />
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
  </div>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_METRIC_UNITS } from '@/constants';

import { convertFileSizeToUnit } from '@/helpers/file/size';
import { saveCsvFile } from '@/helpers/file/files';

import { usePendingHandler } from '@/hooks/query/pending';
import { useInfo } from '@/hooks/store/modules/info';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

import FileDragSelector from '@/components/forms/fields/file-drag-selector.vue';

export default {
  name: 'ImportFileUploadSection',
  components: {
    FileDragSelector,
  },
  props: {
    config: {
      type: Object,
      required: true,
    },
    separator: {
      type: String,
      required: true,
    },
    activeImportFileId: {
      type: String,
      default: null,
    },
    uploading: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:separator', 'choose-file'],
  setup(props, { emit }) {
    const { fileImportMaxSize } = useInfo();
    const { fetchExternalDataTableSchema } = useExternalDataTable();

    const fileImportMaxSizeInKb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.kilobyte)
    ));

    const fileImportMaxSizeInMb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.megabyte)
    ));

    const hasStructure = computed(() => !!props.config?.externalDataTable?.columns?.length);

    const {
      pending: downloading,
      handler: exportTableStructure,
    } = usePendingHandler(async () => {
      const blob = await fetchExternalDataTableSchema({ id: props.config.externalDataTable._id });

      return saveCsvFile(blob, `external-data-table-${props.config.externalDataTable._id}-structure`);
    });

    const updateSeparator = (value) => {
      emit('update:separator', value);
    };

    const chooseFile = (files) => {
      emit('choose-file', files);
    };

    return {
      fileImportMaxSizeInKb,
      fileImportMaxSizeInMb,
      hasStructure,
      downloading,
      exportTableStructure,
      updateSeparator,
      chooseFile,
    };
  },
};
</script>

<style lang="scss" scoped>
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
</style>
