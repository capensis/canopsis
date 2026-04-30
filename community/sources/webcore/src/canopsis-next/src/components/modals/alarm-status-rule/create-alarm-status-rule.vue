<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <alarm-status-rule-form
          v-model="form"
          :disablable="config.disablable"
          :flapping="config.flapping"
          :default-rule="isDefaultRule"
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

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY, DEFAULT_RESOLVE_RULE_ID } from '@/constants';

import { alarmStatusRuleToForm, formToAlarmStatusRule } from '@/helpers/entities/alarm-status-rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAlarmStatusRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { AlarmStatusRuleForm, ModalWrapper, AiChatSidebar },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const { rule, flapping } = props.modal.config;

    const form = ref(alarmStatusRuleToForm(rule, flapping));

    const llmContext = computed(() => (
      config.value.flapping
        ? LLM_SOCKET_CONTEXTS.flappingRule
        : LLM_SOCKET_CONTEXTS.resolveRule
    ));

    const isDefaultRule = computed(() => config.value.rule?._id === DEFAULT_RESOLVE_RULE_ID);

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.rule?._id,
      context: llmContext,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToAlarmStatusRule(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isDefaultRule,
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
