<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <anomaly-monitored-connector-form v-model="form" />
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

import { MODALS } from '@/constants';

import {
  anomalyMonitoredConnectorToForm,
  formToAnomalyMonitoredConnector,
} from '@/helpers/entities/anomaly-monitored-connector/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import AnomalyMonitoredConnectorForm from '@/components/other/anomaly-monitored-connector/form/anomaly-monitored-connector-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAnomalyMonitoredConnector,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AnomalyMonitoredConnectorForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(anomalyMonitoredConnectorToForm(config.value.connector));

    const title = computed(() => (
      config.value.title || t('modals.createAnomalyMonitoredConnector.create.title')
    ));

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToAnomalyMonitoredConnector(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitting,
      isDisabled,

      submit,
      close,
    };
  },
};
</script>
