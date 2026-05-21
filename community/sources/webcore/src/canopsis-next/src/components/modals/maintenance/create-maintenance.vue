<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-4" column>
          <c-alert type="warning">
            {{ config.warningText }}
          </c-alert>
          <maintenance-form v-model="form" />
        </v-layout>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { maintenanceToForm } from '@/helpers/entities/maintenance/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useSubmittableForm } from '@/hooks/submittable-form';

import MaintenanceForm from '@/components/other/maintenance/form/maintenance-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMaintenance,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { MaintenanceForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(maintenanceToForm(config.value.maintenance));

    const title = computed(() => config.value.title ?? t('modals.createMaintenance.setup.title'));

    const submitLabel = computed(() => (
      config.value.maintenance ? t('common.submit') : t('modals.createMaintenance.enableMaintenance')
    ));

    const { submit, submitting, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      title,
      submitLabel,
      submit,
      submitting,
      isDisabled,
      close,
    };
  },
};
</script>
