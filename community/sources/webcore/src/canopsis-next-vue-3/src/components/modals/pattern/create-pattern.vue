<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <pattern-form />
        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
        />
      </template>

      <template #actions>
        <v-btn
          :disabled="submitting"
          variant="text"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting || chatPending"
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
import { computed, toRef } from 'vue';

import {
  MODALS,
  PATTERN_TYPES_TO_LLM_SOCKET_CONTEXTS,
  PATTERN_TYPES_TO_PATTERNS_FIELDS,
} from '@/constants';

import { patternToForm, formToPattern } from '@/helpers/entities/pattern/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import PatternForm from '@/components/forms/pattern-form.vue';

export default {
  name: MODALS.createPattern,

  components: {
    PatternForm,
    AiChatSidebar,
  },

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

    const initialValues = patternToForm(config.value.pattern);

    if (config.value.type) {
      initialValues.type = config.value.type;
    }

    const form = useForm({
      initialValues,
    });
    const { values, handleSubmit } = form;

    const title = computed(() => config.value.title);
    const closeModal = () => modals.hide({ id: props.modal.id });

    const { submit, submitting, submitLabel } = useSubmittableForm({
      item: computed(() => config.value.pattern),
      method: handleSubmit(async () => {
        const result = await config.value.action?.(formToPattern(values.value));

        await config.value.afterSubmit?.(result);

        closeModal();

        return result;
      }),
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form: values,
      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.pattern?._id,
      context: PATTERN_TYPES_TO_LLM_SOCKET_CONTEXTS[config.value.type],
      field: PATTERN_TYPES_TO_PATTERNS_FIELDS[config.value.type],
    });
    const chatPending = computed(() => chatOptions.bind.pending.value);

    return {
      t,
      title,
      submitting,
      closeModal,
      submitLabel,
      submit,
      chatShown,
      chatOptions,
      chatPending,
    };
  },
};
</script>
