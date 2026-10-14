<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <dynamic-info-form
          :rule-id="dynamicInfoId"
          :is-disabled-id-field="isDisabledIdField"
        />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.attrs" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ $t('common.cancel') }}
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

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';
import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';

export default {
  name: MODALS.createDynamicInfo,

  components: {
    AiChatSidebar,
    DynamicInfoForm,
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
    const isNew = computed(() => !config.value.dynamicInfo?._id);
    const title = computed(() => config.value.title || t('modals.createDynamicInfo.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const dynamicInfoId = computed(() => config.value.dynamicInfo?._id);
    const isDisabledIdField = computed(
      () => !!(config.value.isDisabledIdField || config.value.submittingIdField),
    );

    const { values, handleSubmit, submitting } = useForm({
      initialValues: dynamicInfoToForm(config.value.dynamicInfo),
    });

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,
      modal: toRef(props, 'modal'),
      ruleId: dynamicInfoId,
      context: LLM_SOCKET_CONTEXTS.dynamicInfos,
    });
    const submitDisabled = computed(() => submitting.value || chatOptions.value.attrs.pending);

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToDynamicInfo(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();
    });

    useFormConfirmableCloseModal({
      submit,
      form: values,
      close: closeModal,
    });

    return {
      config,
      submitting,
      submitDisabled,
      isNew,
      title,
      submitLabel,
      values,
      dynamicInfoId,
      isDisabledIdField,
      chatShown,
      chatOptions,
      submit,
      closeModal,
    };
  },
};
</script>
