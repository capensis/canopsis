<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title ?? $t('modals.applyEventFilter.title') }}</span>
      </template>
      <template #text="">
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancel"
          />
          <c-event-filter-patterns-field
            v-model="form"
            :excluded-attributes="config.excludedAttributes"
            name="patterns"
            required
            @input="errors.remove('patterns')"
          />
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
import { ref } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { promisedWait } from '@/helpers/async';
import { formGroupsToPatternRules, patternToForm } from '@/helpers/entities/pattern/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.applyEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(patternToForm({ event_pattern: config.value.eventPattern }));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formGroupsToPatternRules(form.value?.groups));
        /**
         * We've added that to avoiding problem with async on the backend side.
         * There is 3000ms timeout on the backend side for sync
         */
        await promisedWait(3000);
        await config.value.afterSubmit?.();
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancel: chatCancel,
    } = useAiChatForm({
      form,
      modalId: props.modal.id,
      ruleId: props.modal.config?.eventsRecordId,
      context: LLM_SOCKET_CONTEXTS.eventRecord,
      field: PATTERNS_FIELDS.event,
    });

    return {
      config,
      form,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancel,

      submit,
      close,
    };
  },
};
</script>
