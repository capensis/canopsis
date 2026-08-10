<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <state-setting-form />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.bind" />
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="closeModal">
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

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import StateSettingForm from '@/components/other/state-setting/form/state-setting-form.vue';
import { LLM_SOCKET_CONTEXTS } from '@/constants';
import { formToStateSetting, stateSettingToForm } from '@/helpers/entities/state-setting/form';
import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

export default {
  components: {
    AiChatSidebar,
    StateSettingForm,
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
    const stateSetting = computed(() => config.value.stateSetting);
    const title = computed(() => config.value.title || t('modals.createStateSetting.create.title'));

    const form = useForm({
      initialValues: stateSettingToForm(stateSetting.value),
    });

    const { values, handleSubmit } = form;

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,
      modal: toRef(props, 'modal'),
      ruleId: stateSetting.value?._id,
      context: LLM_SOCKET_CONTEXTS.stateSettings,
    });
    const chatPending = computed(() => chatOptions.bind.pending.value);

    const closeModal = () => modals.hide({ id: props.modal.id });

    const { submit, submitting, submitLabel } = useSubmittableForm({
      item: stateSetting,
      method: handleSubmit(async () => {
        const result = await config.value.action?.(formToStateSetting(values.value));

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

    return {
      t,
      title,
      submitting,
      chatShown,
      chatOptions,
      chatPending,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
