<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-data-table-form
          v-model="form"
          :is-new="isNew"
          :from-config="fromConfig"
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

import { externalDataTableToForm } from '@/helpers/entities/external-data-table/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ExternalDataTableForm from '@/components/other/external-data-table/form/external-data-table-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createExternalDataTable,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ExternalDataTableForm,
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

    const form = ref(externalDataTableToForm(config.value.externalDataTable));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title || t('modals.createExternalDataTable.create.title'));
    const isNew = computed(() => !config.value.externalDataTable);
    const fromConfig = computed(() => config.value.externalDataTable?.from_config);

    return {
      form,

      isDisabled,
      submitting,

      title,
      isNew,
      fromConfig,

      submit,
      close,
    };
  },
};
</script>
