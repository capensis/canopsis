<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ config.title }}
      </template>
      <template #text="">
        <service-form v-model="form" :prepare-state-setting-form="prepareStateSettingForm" />
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
import { ref } from 'vue';

import { ENTITY_TYPES, MODALS, VALIDATION_DELAY } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/entities/service/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useInnerModal } from '@/hooks/modals';

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

    useFormConfirmableCloseModal({ form, submit, close });

    const prepareStateSettingForm = service => ({
      ...formToService(service),
      type: ENTITY_TYPES.service,
      _id: service._id,
    });

    return {
      config,
      form,
      isDisabled,
      submitting,

      close,
      prepareStateSettingForm,
      submit,
    };
  },
};
</script>
