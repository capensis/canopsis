<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <dynamic-info-form
          v-model="form"
          :rule-id="dynamicInfoId"
          :is-disabled-id-field="submittingIdField"
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
          :disabled="chatOptions.bind.pending"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    AiChatSidebar,
    DynamicInfoForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(dynamicInfoToForm(config.value.dynamicInfo));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.dynamicInfo?._id,
      context: LLM_SOCKET_CONTEXTS.dynamicInfos,
    });

    const dynamicInfoId = computed(() => config.value.dynamicInfo?._id);
    const title = computed(() => config.value.title || t('modals.createDynamicInfo.create.title'));
    const submittingIdField = computed(() => config.value.submittingIdField);

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.dynamicInfo,
      method: async () => {
        const result = await config.value.action?.(formToDynamicInfo(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      dynamicInfoId,
      title,
      submittingIdField,
      submitting,
      chatShown,
      chatOptions,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>
