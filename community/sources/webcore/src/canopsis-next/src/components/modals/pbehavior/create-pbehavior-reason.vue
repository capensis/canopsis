<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createPbehaviorReason.title') }}</span>
      </template>
      <template #text="">
        <pbehavior-reason-form v-model="form" />
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

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorReasonToForm, formToPbehaviorReason } from '@/helpers/entities/pbehavior/reason/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PbehaviorReasonForm from '@/components/other/pbehavior/reasons/form/pbehavior-reason-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorReason,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PbehaviorReasonForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(pbehaviorReasonToForm(config.value.pbehaviorReason));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToPbehaviorReason(form.value));
        }

        close();
      },
    });

    return {
      form,
      isDisabled,
      submitting,
      close,
      submit,
    };
  },
};
</script>
