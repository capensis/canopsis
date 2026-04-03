<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancelPending"
          />
          <link-rule-form
            v-model="form"
            :rule-id="config.linkRule?._id"
          />
        </div>
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

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { linkRuleToForm, formToLinkRule } from '@/helpers/entities/link/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import LinkRuleForm from '@/components/other/link-rule/form/link-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createLinkRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    LinkRuleForm,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.linkRule;

    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(linkRuleToForm(config.value.linkRule));

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
      ruleId: props.modal.config?.linkRule?._id,
      context: LLM_SOCKET_CONTEXTS.linkRule,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const linkRule = await config.value.action?.(formToLinkRule(form.value));

        close();

        return linkRule;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });
    useEntityInfoPropertyFetching();

    const title = computed(() => config.value.title ?? t('modals.createLinkRule.create.title'));

    return {
      type,

      config,

      form,

      isDisabled,
      submitting,

      title,

      chatPending,
      chatPendingTexts,
      chatCancelPending,

      submit,
      close,
    };
  },
};
</script>
