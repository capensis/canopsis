<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <template-testing-test-variables-wrapper
          v-model="form"
          :rule-id="dynamicInfoId"
          :type="type"
        >
          <template #default="{ templateVars, copyVars }">
            <dynamic-info-form
              v-model="form"
              :is-disabled-id-field="isDisabledIdField"
              :template-vars="templateVars"
              :copy-vars="copyVars"
            />
            <ai-chat-sidebar
              v-if="chatShown"
              v-bind="chatOptions.bind"
              v-on="chatOptions.on"
            />
          </template>
        </template-testing-test-variables-wrapper>
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
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DynamicInfoForm,
    AiChatSidebar,
    TemplateTestingTestVariablesWrapper,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.dynamicInfo;

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
    const isDisabledIdField = computed(() => config.value.isDisabledIdField);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
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
      config,
      dynamicInfoId,
      type,
      title,
      isDisabledIdField,
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
