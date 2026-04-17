<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <tag-form
          v-model="form"
          :is-imported="isImported"
          :is-new="isNew"
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
import { ref, computed, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { tagToForm, formToTag } from '@/helpers/entities/tag/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import TagForm from '@/components/other/tag/form/tag-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTag,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { TagForm, ModalWrapper, AiChatSidebar },
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

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
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
        const result = await config.value.action?.(formToTag(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
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
      chatShown,
      chatOptions,
      submit,
      close,
    };
  },
};
</script>
