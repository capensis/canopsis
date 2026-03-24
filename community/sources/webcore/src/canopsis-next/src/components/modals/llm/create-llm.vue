<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <llm-form v-model="form" :is-new="isNew" />
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
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { llmToForm, formToLlm } from '@/helpers/entities/llm/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import LlmForm from '@/components/other/llm/form/llm-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createLlm,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    LlmForm,
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

    const form = ref(llmToForm(config.value.llm));

    const isNew = computed(() => !config.value.llm?._id);

    const title = computed(() => (
      config.value.title || (
        isNew.value
          ? t('modals.createLlm.create.title')
          : t('modals.createLlm.edit.title')
      )
    ));

    const submitLabel = computed(() => (
      isNew.value ? t('common.add') : t('common.submit')
    ));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToLlm(form.value));

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

      submit,
      close,
    };
  },
};
</script>
