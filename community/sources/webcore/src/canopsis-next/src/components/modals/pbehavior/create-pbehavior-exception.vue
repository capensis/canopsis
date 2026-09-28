<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-exception-form v-model="form" />
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

import { formToPbehaviorException, pbehaviorExceptionToForm } from '@/helpers/entities/pbehavior/exception/form';

import { useI18n } from '@/hooks/i18n';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PbehaviorExceptionForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorException,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PbehaviorExceptionForm,
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

    const form = ref(pbehaviorExceptionToForm(config.value.pbehaviorException));

    const isNew = computed(() => !config.value.pbehaviorException?._id);

    const title = computed(() => (
      config.value.title || t(isNew.value
        ? 'modals.createPbehaviorException.create.title'
        : 'modals.createPbehaviorException.edit.title')
    ));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.pbehaviorException,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToPbehaviorException(form.value));
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
