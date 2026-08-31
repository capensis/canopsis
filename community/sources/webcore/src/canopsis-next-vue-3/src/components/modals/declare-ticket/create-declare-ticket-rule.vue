<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <declare-ticket-rule-form :rule-id="config.declareTicketRule?._id" />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.attrs" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitDisabled"
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

import { LLM_SOCKET_CONTEXTS, MODALS } from '@/constants';

import {
  declareTicketRuleToForm,
  formToDeclareTicketRule,
} from '@/helpers/entities/declare-ticket/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';
import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';

export default {
  name: MODALS.createDeclareTicketRule,

  components: {
    DeclareTicketRuleForm,
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
    const close = () => modals.hide({ id: props.modal.id });

    const form = useForm({
      initialValues: declareTicketRuleToForm(config.value.declareTicketRule),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.declareTicketRule?._id,
      context: LLM_SOCKET_CONTEXTS.declareTicketRule,
    });

    const isNew = computed(() => !config.value.declareTicketRule?._id);
    const title = computed(
      () => config.value.title ?? t('modals.createDeclareTicketRule.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToDeclareTicketRule(values.value));

      await config.value.afterSubmit?.(result);

      close();

      return result;
    });
    const submitDisabled = computed(() => submitting.value || chatOptions.value.attrs.pending);

    useFormConfirmableCloseModal({ form: values, submit, close });

    return {
      t,
      config,
      isNew,
      title,
      submitting,
      submitDisabled,
      chatShown,
      chatOptions,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>
