<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancel"
          />
          <tag-form
            v-model="form"
            :is-imported="isImported"
            :is-new="isNew"
          />
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
import { ref, computed } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { tagToForm, formToTag } from '@/helpers/entities/tag/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import TagForm from '@/components/other/tag/form/tag-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTag,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { TagForm, ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(tagToForm(config.value.tag));

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancel: chatCancel,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: props.modal.id,
      ruleId: props.modal.config?.tag?._id,
      context: LLM_SOCKET_CONTEXTS.alarmTag,
    });

    const isNew = computed(() => !config.value.tag?._id);
    const title = computed(() => (
      config.value.title || t('modals.createTag.create.title')
    ));

    const isImported = computed(() => config.value.isImported);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToTag(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      isNew,
      isImported,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancel,
      submit,
      close,
    };
  },
};
</script>
