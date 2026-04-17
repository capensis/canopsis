<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <remediation-instruction-approval-alert
          v-if="hasApproval && isChangesByCurrentUser"
          :user-name="alertUserName"
          :comment="alertComment"
          :dismissed="isChangesDismissed"
          class="mb-3"
        />
        <remediation-instruction-form
          v-model="form"
          :disabled="disabled"
          :is-new="isNew"
          :required-approve="requiredInstructionApprove"
          :rule-id="config.remediationInstruction?._id"
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
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled || chatOptions.bind.pending"
          :loading="submitting"
          class="primary"
          type="submit"
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, toRef } from 'vue';
import { createNamespacedHelpers } from 'vuex';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import {
  formToRemediationInstructionRequest,
  remediationInstructionErrorsToForm,
  remediationInstructionToFullForm,
} from '@/helpers/entities/remediation/instruction/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useAuth } from '@/hooks/auth';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';
import RemediationInstructionApprovalAlert from '@/components/other/remediation/instructions/partials/approval-alert.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  name: MODALS.createRemediationInstruction,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
    AiChatSidebar,
    RemediationInstructionForm,
    RemediationInstructionApprovalAlert,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.instruction;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const { currentUser } = useAuth();

    const form = ref(remediationInstructionToFullForm(config.value.remediationInstruction));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.remediationInstruction?._id,
      context: LLM_SOCKET_CONTEXTS.instruction,
    });

    const title = computed(() => config.value.title || t('modals.createRemediationInstruction.create.title'));
    const disabled = computed(() => config.value.disabled);
    const isNew = computed(() => !config.value.remediationInstruction?._id);
    const approval = computed(() => config.value.remediationInstruction?.approval);
    const hasApproval = computed(() => !!approval.value);
    const isChangesDismissed = computed(() => !!approval.value?.dismissed_by);
    const isChangesByCurrentUser = computed(() => approval.value?.requested_by?._id === currentUser.value._id);

    const alertUserName = computed(() => {
      const { dismissed_by: dismissedBy, requested_by: requestedBy } = approval.value ?? {};
      return dismissedBy?.display_name ?? requestedBy?.display_name;
    });

    const alertComment = computed(() => approval.value?.dismiss_comment ?? approval.value?.comment);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToRemediationInstructionRequest(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
      errorsToValidation: err => remediationInstructionErrorsToForm(err, form.value),
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      disabled,
      approval,
      hasApproval,
      isChangesDismissed,
      isChangesByCurrentUser,
      alertUserName,
      alertComment,
      isDisabled,
      submitting,
      chatShown,
      chatOptions,
      submit,
    };
  },
  computed: {
    ...mapGetters({
      requiredInstructionApprove: 'requiredInstructionApprove',
    }),
  },
};
</script>
