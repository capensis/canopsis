<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ config.title ?? $t('modals.applyEventFilter.title') }}</span>
      </template>
      <template #text="">
        <div class="position-relative">
          <c-progress-overlay :pending="eventPatternAttributesPending" />
          <v-alert v-if="config.infoAlert" class="mb-4" type="info">
            {{ config.infoAlert }}
          </v-alert>
          <c-event-filter-patterns-field
            v-model="form"
            :attributes="eventPatternAttributes"
            :required="config.required ?? true"
            name="patterns"
            @input="errors.remove('patterns')"
          />
        </div>
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
import { ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { promisedWait } from '@/helpers/async';
import { formGroupsToPatternRules, patternToForm } from '@/helpers/entities/pattern/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.applyEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper, AiChatSidebar },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(patternToForm({ event_pattern: config.value.eventPattern }));

    const { fetchEventRecordPatternFields } = usePatternsFields();

    const {
      pending: eventPatternAttributesPending,
      eventPatternAttributes,
    } = usePatternsFieldsFetching(fetchEventRecordPatternFields);

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.eventsRecordId,
      context: LLM_SOCKET_CONTEXTS.eventRecord,
      field: PATTERNS_FIELDS.event,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formGroupsToPatternRules(form.value?.groups));

        /**
         * We've added that to avoiding problem with async on the backend side.
         * There is 3000ms timeout on the backend side for sync
         */
        await promisedWait(3000);
        await config.value.afterSubmit?.(result);
        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      config,
      form,
      isDisabled,
      submitting,
      chatShown,
      chatOptions,

      eventPatternAttributes,
      eventPatternAttributesPending,

      submit,
      close,
    };
  },
};
</script>
