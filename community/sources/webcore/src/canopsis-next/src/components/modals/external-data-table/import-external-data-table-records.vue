<template>
  <v-form class="import-external-data-table-records-form" @submit.prevent="submitButton.action">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-data-table-general-info-form :form="config.externalDataTable" />
        <import-file-upload-section
          :separator="separator"
          :config="config"
          :active-import-file-id="activeImportFileId"
          :uploading="uploading"
          @update:separator="separator = $event"
          @choose-file="chooseFile"
        />
        <import-preview-section
          v-if="activeImportFileId"
          v-model="form"
          :records="records"
          :pending="pending"
          :options="options"
          :total-items="meta.total_count"
          :separator="separator"
          @update:options="updateOptions"
        />
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
          :disabled="!activeImportFileId || submitButton.loading || isDisabled"
          class="primary"
          type="submit"
        >
          {{ submitButton.text }}
          <v-progress-circular
            v-if="submitButton.loading"
            class="ml-2"
            size="24"
            indeterminate
          />
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import {
  useExternalDataTableImportFile,
} from '@/components/other/external-data-table/hooks/external-data-table-import-file';

import ExternalDataTableGeneralInfoForm
  from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

import ImportFileUploadSection from './partials/import-file-upload-section.vue';
import ImportPreviewSection from './partials/import-preview-section.vue';

export default {
  name: MODALS.importExternalDataTableRecords,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ExternalDataTableGeneralInfoForm,
    ImportFileUploadSection,
    ImportPreviewSection,
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

    const {
      separator,
      activeImportFileId,
      form,
      needPreview,
      isReadyForComplete,
      records,
      meta,
      pending,
      options,
      updatingPreview,
      uploading,
      updateOptions,
      updatePreview,
      chooseFile,
      completeImport,
      toggleOnNeedPreview,
    } = useExternalDataTableImportFile({ config });

    const { submit: complete, isDisabled, submitting } = useSubmittableForm({
      form: separator,
      method: async () => {
        await completeImport(config.value.afterSubmit);
        close();
      },
    });

    const title = computed(() => config.value.title || t('modals.createExternalDataTableRecord.create.title'));

    const submitButton = computed(() => {
      if (updatingPreview.value) {
        return {
          text: t('externalData.loadingPreview'),
          loading: true,
        };
      }

      if (needPreview.value) {
        return {
          text: t('externalData.updatePreview'),
          action: updatePreview,
        };
      }

      return {
        text: t('common.import'),
        loading: submitting.value,
        disabled: !isReadyForComplete.value,
        action: complete,
      };
    });

    return {
      config,
      separator,
      activeImportFileId,
      form,
      isDisabled,
      submitting,
      pending,
      meta,
      records,
      options,
      updateOptions,
      title,
      submitButton,
      updatePreview,
      uploading,
      close,
      chooseFile,
      toggleOnNeedPreview,
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
}
</style>
