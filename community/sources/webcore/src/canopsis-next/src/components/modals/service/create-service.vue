<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ config.title }}
      </template>
      <template #text="">
        <c-progress-overlay :pending="pending" />
        <service-form
          v-model="form"
          :prepare-state-setting-form="prepareStateSettingForm"
          :template-vars="templateVars"
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
import { ref, onMounted, toRef } from 'vue';

import { ENTITY_TYPES, LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/entities/service/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useTemplateVarsList } from '@/hooks/vars/template';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import ServiceForm from '@/components/other/service/form/service-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createService,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ServiceForm, ModalWrapper, AiChatSidebar },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(serviceToForm(config.value.item));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.item?._id,
      context: LLM_SOCKET_CONTEXTS.entityService,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToService(form.value));

        await config.value.afterSubmit?.(result);

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
      chatShown,
      chatOptions,

      close,
      prepareStateSettingForm,
      submit,
    };
  },
};
</script>
