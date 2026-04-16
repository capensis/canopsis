<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <pattern-form v-model="form" />
        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
          v-on="chatOptions.on"
        />
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
          :disabled="isDisabled || chatOptions.bind.pending"
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
import { ref, toRef } from 'vue';

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

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import PatternForm from '@/components/forms/pattern-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPattern,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PatternForm,
    AiChatSidebar,
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
        const result = await config.value.action?.(formToPattern(form.value));

        await config.value.afterSubmit?.(result);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
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
      chatShown,
      chatOptions,
    };
  },
};
</script>
