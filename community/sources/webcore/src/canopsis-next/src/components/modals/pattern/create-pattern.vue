<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancel"
          />
          <pattern-form v-model="form" />
        </div>
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
          :disabled="isDisabled || chatPending"
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

import {
  MODALS,
  VALIDATION_DELAY,
  PATTERN_TYPES_TO_LLM_SOCKET_CONTEXTS,
  PATTERN_TYPES_TO_PATTERNS_FIELDS,
} from '@/constants';

import { patternToForm, formToPattern } from '@/helpers/entities/pattern/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useAiChatForm } from '@/hooks/ai/ai-chat-form';

import PatternForm from '@/components/forms/pattern-form.vue';
import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPattern,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PatternForm,
    PatternProgress,
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

    const form = ref(patternToForm(config.value.pattern));

    if (config.value.type) {
      form.value.type = config.value.type;
    }

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToPattern(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancel: chatCancel,
    } = useAiChatForm({
      form,
      modalId: props.modal.id,
      ruleId: props.modal.config?.pattern?._id,
      context: PATTERN_TYPES_TO_LLM_SOCKET_CONTEXTS[config.value.type],
      field: PATTERN_TYPES_TO_PATTERNS_FIELDS[config.value.type],
    });

    return {
      config,
      form,
      isDisabled,
      submitting,
      close,
      submit,
      chatPending,
      chatPendingTexts,
      chatCancel,
    };
  },
};
</script>
