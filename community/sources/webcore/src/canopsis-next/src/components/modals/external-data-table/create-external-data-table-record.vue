<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-data-table-record-form
          v-model="form"
          :external-data-table="config.externalDataTable"
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

import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  externalDataTableRecordToForm,
  formToExternalDataTableRecord,
} from '@/helpers/entities/external-data-table/record/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ExternalDataTableRecordForm from '@/components/other/external-data-table/form/external-data-table-record-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createExternalDataTableRecord,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ExternalDataTableRecordForm,
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

    const columnConfigs = computed(() => config.value.externalDataTable.column_configs ?? []);

    const form = ref(externalDataTableRecordToForm(
      config.value.externalDataTableRecord,
      columnConfigs.value,
    ));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToExternalDataTableRecord(form.value, columnConfigs.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title || t('modals.createExternalDataTableRecord.create.title'));

    return {
      config,

      form,

      isDisabled,
      submitting,

      title,

      submit,
      close,
    };
  },
};
</script>
