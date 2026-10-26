<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <meta-alarm-rule-form :rule-id="ruleId" />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.attrs" />
      </template>

      <template #actions>
        <v-btn
          :disabled="isSubmitting"
          variant="text"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitDisabled"
          :loading="isSubmitting"
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

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

export default {
  name: MODALS.createMetaAlarmRule,

  components: {
    AiChatSidebar,
    MetaAlarmRuleForm,
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
    const rule = computed(() => config.value.rule);
    const ruleId = computed(() => rule.value?._id);
    const isNew = computed(() => !ruleId.value || config.value.isDuplicating);
    const title = computed(() => config.value.title ?? t('modals.metaAlarmRule.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));
    const form = useForm({
      initialValues: metaAlarmRuleToForm(rule.value),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,
      modal: toRef(props, 'modal'),
      ruleId: ruleId.value,
      context: LLM_SOCKET_CONTEXTS.metaAlarmRule,
    });
    const submitDisabled = computed(() => isSubmitting.value || chatOptions.value.attrs.pending);

    const closeModal = () => modals.hide({ id: props.modal.id });
    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToMetaAlarmRule(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      ruleId,
      submitLabel,
      isSubmitting,
      submitDisabled,
      chatShown,
      chatOptions,
      submit,
      closeModal,
    };
  },
};
</script>
