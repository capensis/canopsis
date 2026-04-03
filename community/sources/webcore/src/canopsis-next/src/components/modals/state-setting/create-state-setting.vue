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
          <state-setting-form v-model="form" />
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

import { stateSettingToForm, formToStateSetting } from '@/helpers/entities/state-setting/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import StateSettingForm from '@/components/other/state-setting/form/state-setting-form.vue';
import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createStateSetting,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    StateSettingForm,
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

    const form = ref(stateSettingToForm(config.value.stateSetting));

    const title = computed(() => config.value.title || t('modals.createStateSetting.create.title'));

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancelPending: chatCancelPending,
    } = useAiChatForm({
      form,
      modalId: props.modal.id,
      ruleId: props.modal.config?.stateSetting?._id,
      context: LLM_SOCKET_CONTEXTS.stateSettings,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToStateSetting(form.value));
        close();
      },
    });

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
