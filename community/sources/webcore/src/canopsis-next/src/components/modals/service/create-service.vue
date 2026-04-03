<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ config.title }}
      </template>
      <template #text="">
        <c-progress-overlay :pending="pending" />
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancelPending"
          />
          <service-form
            v-model="form"
            :prepare-state-setting-form="prepareStateSettingForm"
            :template-vars="templateVars"
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
import { computed, ref, onMounted } from 'vue';

import { ENTITY_TYPES, LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/entities/service/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useTemplateVarsList } from '@/hooks/vars/template';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import ServiceForm from '@/components/other/service/form/service-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createService,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ServiceForm, ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(serviceToForm(config.value.item));

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancelPending: chatCancelPending,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: props.modal.id,
      ruleId: props.modal.config?.item?._id,
      context: LLM_SOCKET_CONTEXTS.entityService,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToService(form.value));
        close();
      },
    });

    const { vars: templateVars, pending, fetchList } = useTemplateVarsList({
      type: ENTITY_TYPES.service,
    });

    useFormConfirmableCloseModal({ form, submit, close });

    /**
     * Prepares a service object for state setting form by transforming it to the proper format
     *
     * @param {Object} service - The service object to prepare
     * @param {string} service._id - The service identifier
     * @returns {Object} The prepared service object with form data, entity type, and ID
     */
    const prepareStateSettingForm = service => ({
      ...formToService(service),
      type: ENTITY_TYPES.service,
      _id: service._id,
    });

    onMounted(fetchList);

    return {
      config,
      form,
      templateVars,
      pending,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancelPending,

      close,
      prepareStateSettingForm,
      submit,
    };
  },
};
</script>
