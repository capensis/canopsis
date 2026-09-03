<template>
  <v-form class="import-external-data-table-records-form" @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
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
          :form="form"
          :original-form="originalForm"
          :records="records"
          :pending="pending"
          :options="options"
          :total-items="meta.total_count"
          :separator="separator"
          @update:form="updateForm"
          @input="validate"
          @update:options="updateOptions"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="!activeImportFileId || submitButton.disabled || submitButton.loading || submitting"
          :loading="submitButton.loading"
          color="primary"
          type="submit"
        >
          {{ submitButton.text }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import { useExternalDataTableImportFile } from '@/components/other/external-data-table/hooks/external-data-table-import-file';

import ExternalDataTableGeneralInfoForm from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';
import ImportFileUploadSection from '@/components/other/external-data-table/partials/import-file-upload-section.vue';
import ImportPreviewSection from '@/components/other/external-data-table/partials/import-preview-section.vue';

export default {
  components: {
    ExternalDataTableGeneralInfoForm,
    ImportFileUploadSection,
    ImportPreviewSection,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});

    const {
      formContext,
      separator,
      activeImportFileId,
      form,
      originalForm,
      needPreview,
      hasErrors,
      records,
      meta,
      pending,
      options,
      uploading,
      updateOptions,
      updatePreview,
      chooseFile,
      completeImport,
      validate,
    } = useExternalDataTableImportFile({ config });

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submitting = computed(() => pending.value || uploading.value);

    const submit = async () => {
      const { valid } = await formContext.validate();

      if (valid) {
        if (needPreview.value) {
          await updatePreview();

          return;
        }

        await completeImport(config.value.afterSubmit);

        closeModal();
      }
    };

    const title = computed(
      () => config.value.title || t('modals.importExternalDataTableRecords.title'),
    );

    const submitButton = computed(() => {
      if (needPreview.value) {
        if (submitting.value) {
          return {
            text: t('externalData.loadingPreview'),
            loading: true,
          };
        }

        return {
          text: t('externalData.updatePreview'),
          disabled: hasErrors.value,
        };
      }

      return {
        text: t('common.import'),
        loading: submitting.value,
      };
    });

    const updateForm = (newForm) => {
      form.value = newForm;
    };

    return {
      t,
      config,
      separator,
      activeImportFileId,
      form,
      originalForm,
      submit,
      submitting,
      pending,
      meta,
      records,
      options,
      updateOptions,
      title,
      submitButton,
      uploading,
      closeModal,
      chooseFile,
      validate,
      updateForm,
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
