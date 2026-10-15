<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <external-data-table-record-form :external-data-table="config.externalDataTable" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  externalDataTableRecordToForm,
  formToExternalDataTableRecord,
} from '@/helpers/entities/external-data-table/record/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import ExternalDataTableRecordForm from '@/components/other/external-data-table/form/external-data-table-record-form.vue';

export default {
  components: {
    ExternalDataTableRecordForm,
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
    const columnConfigs = computed(() => config.value.externalDataTable.column_configs ?? []);
    const externalDataTableRecord = computed(() => config.value.externalDataTableRecord);
    const isNew = computed(() => !externalDataTableRecord.value?._id);
    const title = computed(
      () => config.value.title || t('modals.createExternalDataTableRecord.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: externalDataTableRecordToForm(
        externalDataTableRecord.value,
        columnConfigs.value,
      ),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(
        formToExternalDataTableRecord(values.value, columnConfigs.value),
      );

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      config,
      submitting,
      title,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
