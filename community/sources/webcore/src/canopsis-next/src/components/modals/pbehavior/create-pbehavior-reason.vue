<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-reason-form v-model="form" />
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
          :disabled="submitting"
          :loading="submitting"
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
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorReasonToForm, formToPbehaviorReason } from '@/helpers/entities/pbehavior/reason/form';

import { useI18n } from '@/hooks/i18n';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
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
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(pbehaviorReasonToForm(config.value.pbehaviorReason));

    const isNew = computed(() => !config.value.pbehaviorReason?._id);

    const title = computed(() => (
      config.value.title || t(isNew.value
        ? 'modals.createPbehaviorReason.create.title'
        : 'modals.createPbehaviorReason.edit.title')
    ));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.pbehaviorReason,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToPbehaviorReason(form.value));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitting,
      close,
      submitLabel,
      submit,
    };
  },
};
</script>
