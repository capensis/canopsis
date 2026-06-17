<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <idle-rule-form v-model="form" />
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
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { formToIdleRule, idleRuleToForm } from '@/helpers/entities/idle-rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createIdleRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    IdleRuleForm,
    ModalWrapper,
    AiChatSidebar,
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

    const form = ref(idleRuleToForm(config.value.idleRule));

    const title = computed(() => config.value.title || t('modals.createAlarmIdleRule.create.title'));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.idleRule?._id,
      context: LLM_SOCKET_CONTEXTS.idleRule,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToIdleRule(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useEntityInfoPropertyFetching();
    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
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
