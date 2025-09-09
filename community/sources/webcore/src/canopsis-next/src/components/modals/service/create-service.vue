<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ config.title }}
      </template>
      <template #text="">
        <c-progress-overlay :pending="pending" />
        <service-form
          v-model="form"
          :prepare-state-setting-form="prepareStateSettingForm"
          :template-vars="templateVars"
        />
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
import { ref, onMounted } from 'vue';

import { ENTITY_TYPES, MODALS, VALIDATION_DELAY } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/entities/service/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useInnerModal } from '@/hooks/modals';

import { useTemplateVarsList } from '@/components/other/template-testing/hooks/template-test-variables-wrapper';

import ServiceForm from '@/components/other/service/form/service-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createService,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ServiceForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(serviceToForm(config.value.item));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToService(form.value));
        close();
      },
    });

    const { templateVars, pending, fetchList } = useTemplateVarsList({
      type: ENTITY_TYPES.service,
    });

    useFormConfirmableCloseModal({ form, submit, close });

    /**
     * Prepares a service object for state setting form by transforming it to the proper format
     *
     * @param {Object} service - The service object to prepare
     * @param {string} service._id - The service identifier
     * @returns {Object} The prepared service object with form data, entity type, and ID
     */
    const prepareStateSettingForm = service => ({
      ...formToService(service),
      type: ENTITY_TYPES.service,
      _id: service._id,
    });

    onMounted(fetchList);

    return {
      config,
      form,
      templateVars,
      pending,
      isDisabled,
      submitting,

      close,
      prepareStateSettingForm,
      submit,
    };
  },
};
</script>
