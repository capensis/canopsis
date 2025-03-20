<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-data-table-general-info-form :form="config.externalDataTable" />
        <c-csv-separator-field v-model="separator" />
        <file-drag-selector
          v-bind="$attrs"
          :file-type-label="$t('common.fileSelector.fileTypes.csv')"
          :max-file-size="fileImportMaxSizeInKb"
          accept=".csv"
          required
          @change="chooseFile"
        >
          <template #label="{ on }">
            <v-layout column>
              <span class="text-subtitle-2">
                {{ $t('common.fileSelector.dragAndDrop.label') }}
                <a v-on="on">
                  {{ $t('common.fileSelector.dragAndDrop.labelAction') }}
                </a>
                {{ $t('common.fileSelector.fileTypes.csv') }}
                ({{ $t('common.fileSelector.fileSizeMb', { size: fileImportMaxSizeInMb }) }})
              </span>
              <span class="text-subtitle-2">{{ $t('externalData.importFileDescription') }}</span>
            </v-layout>
          </template>
        </file-drag-selector>
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
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { CSV_SEPARATORS, EXTERNAL_METRIC_UNITS, MODALS, VALIDATION_DELAY } from '@/constants';

import { convertFileSizeToUnit } from '@/helpers/file/size';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useInfo } from '@/hooks/store/modules/info';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

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

    const { fileImportMaxSize } = useInfo();
    const {
      createExternalDataTableImport,
      fetchExternalDataTableImportData,
    } = useExternalDataTable();

    const separator = ref(CSV_SEPARATORS.comma);

    const fileImportMaxSizeInKb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.kilobyte)
    ));

    const fileImportMaxSizeInMb = computed(() => (
      convertFileSizeToUnit(fileImportMaxSize.value, EXTERNAL_METRIC_UNITS.megabyte)
    ));

    /* const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close }); */

    const title = computed(() => config.value.title || t('modals.createExternalDataTableRecord.create.title'));

    const chooseFile = async ([file] = []) => {
      const data = {
        separator: separator.value,
        file,
      };

      const { _id: id } = await createExternalDataTableImport({ id: config.value.externalDataTable._id, data });
      fetchExternalDataTableImportData({ id });
    };

    return {
      config,
      fileImportMaxSize,
      fileImportMaxSizeInKb,
      fileImportMaxSizeInMb,
      separator,

      // form,

      // isDisabled,
      // submitting,

      title,

      // submit,
      close,
      chooseFile,
    };
  },
};
</script>
