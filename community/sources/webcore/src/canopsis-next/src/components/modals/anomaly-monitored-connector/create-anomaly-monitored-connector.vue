<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <anomaly-monitored-connector-form />
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
  anomalyMonitoredConnectorToForm,
  formToAnomalyMonitoredConnector,
} from '@/helpers/entities/anomaly-monitored-connector/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AnomalyMonitoredConnectorForm from '@/components/other/anomaly-monitored-connector/form/anomaly-monitored-connector-form.vue';

export default {
  components: { AnomalyMonitoredConnectorForm },

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
    const connector = computed(() => config.value.connector);
    const isNew = computed(() => !connector.value);
    const title = computed(
      () => config.value.title || t('modals.createAnomalyMonitoredConnector.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: anomalyMonitoredConnectorToForm(connector.value),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToAnomalyMonitoredConnector(values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      submitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
