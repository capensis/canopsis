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
            @cancel="chatCancelPending"
          />
          <idle-rule-form v-model="form" />
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
import { computed, ref } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { formToIdleRule, idleRuleToForm } from '@/helpers/entities/idle-rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createIdleRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    IdleRuleForm,
    ModalWrapper,
    PatternProgress,
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

    const form = ref(idleRuleToForm(config.value.idleRule));

    const title = computed(() => config.value.title || t('modals.createAlarmIdleRule.create.title'));

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancelPending: chatCancelPending,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: props.modal.id,
      ruleId: props.modal.config?.idleRule?._id,
      context: LLM_SOCKET_CONTEXTS.idleRule,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToIdleRule(form.value));
        close();
      },
    });

    useEntityInfoPropertyFetching();
    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancelPending,
      submit,
      close,
    };
  },
};
</script>
