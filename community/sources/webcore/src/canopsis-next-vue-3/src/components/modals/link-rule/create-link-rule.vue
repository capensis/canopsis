<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <link-rule-form :rule-id="config.linkRule?._id" />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.bind" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting || chatOptions.bind.pending"
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

import LinkRuleForm from '@/components/other/link-rule/form/link-rule-form.vue';
import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import { LLM_SOCKET_CONTEXTS } from '@/constants';
import { linkRuleToForm, formToLinkRule, isDefaultSourceCode } from '@/helpers/entities/link/form';
import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

export default {
  components: {
    LinkRuleForm,
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
    const linkRule = computed(() => config.value.linkRule);
    const isNew = computed(() => !linkRule.value?._id);
    const title = computed(() => config.value.title ?? t('modals.createLinkRule.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: linkRuleToForm(linkRule.value),
    });
    const { values, handleSubmit, submitting, setFieldError } = form;

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.linkRule?._id,
      context: LLM_SOCKET_CONTEXTS.linkRule,
    });

    const closeModal = () => modals.hide({ id: props.modal.id });

    const validateLinks = () => {
      const hasLinks = !!values.value.links?.length;
      const hasCustomSourceCode = !isDefaultSourceCode(values.value.source_code);

      if (hasLinks || hasCustomSourceCode) {
        setFieldError('links', undefined);

        return true;
      }

      setFieldError('links', t('linkRule.linksEmptyError'));

      return false;
    };

    const submit = handleSubmit(async () => {
      if (!validateLinks()) {
        return undefined;
      }

      const result = await config.value.action?.(formToLinkRule(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });
    useEntityInfoPropertyFetching();

    return {
      t,
      config,
      title,
      submitLabel,
      submitting,
      chatShown,
      chatOptions,
      closeModal,
      submit,
    };
  },
};
</script>
