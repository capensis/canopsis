<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <declare-ticket-rule-form
          v-model="form"
          :rule-id="config.declareTicketRule?._id"
        />
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
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { declareTicketRuleToForm, formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDeclareTicketRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DeclareTicketRuleForm,
    ModalWrapper,
    AiChatSidebar,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(declareTicketRuleToForm(config.value.declareTicketRule));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.declareTicketRule?._id,
      context: LLM_SOCKET_CONTEXTS.declareTicketRule,
    });

    const isNew = computed(() => !config.value.declareTicketRule?._id);
    const title = computed(() => config.value.title ?? t('modals.createDeclareTicketRule.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToDeclareTicketRule(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      isDisabled,
      submitting,
      chatShown,
      chatOptions,
      submit,
      close,
    };
  },
};
</script>
