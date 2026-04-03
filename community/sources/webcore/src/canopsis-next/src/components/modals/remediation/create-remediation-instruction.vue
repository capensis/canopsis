<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
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
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancelPending"
          />
          <remediation-instruction-form
            v-model="form"
            :disabled="disabled"
            :is-new="isNew"
            :required-approve="requiredInstructionApprove"
            :rule-id="config.remediationInstruction?._id"
          />
        </div>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled || chatPending"
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
import { computed, ref } from 'vue';
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

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
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
    PatternProgress,
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
        const data = await config.value.action?.(formToRemediationInstructionRequest(form.value));

        close();

        return data;
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
      chatPending,
      chatPendingTexts,
      chatCancelPending,
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
