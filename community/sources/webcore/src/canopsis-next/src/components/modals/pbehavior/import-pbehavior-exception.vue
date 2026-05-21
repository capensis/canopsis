<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-exception-import-form v-model="form" />
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
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorExceptionImportToForm } from '@/helpers/entities/pbehavior/exception/form';

import { useI18n } from '@/hooks/i18n';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PbehaviorExceptionImportForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-import-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importPbehaviorException,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PbehaviorExceptionImportForm,
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

    const form = ref(pbehaviorExceptionImportToForm(config.value.pbehaviorException));

    const title = computed(() => config.value.title || t('modals.importPbehaviorException.title'));

    const submitLabel = computed(() => config.value.submitLabel || t('common.import'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitLabel,
      isDisabled,
      submitting,
      close,
      submit,
    };
  },
};
</script>
