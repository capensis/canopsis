<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <meta-alarm-rule-form
          v-model="form"
          :rule-id="ruleId"
          :disabled-id-field="!isNew"
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
          :disabled="submitting || chatOptions.bind.pending"
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
import { computed, ref, onMounted, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfos } from '@/hooks/store/modules/entity-infos';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    MetaAlarmRuleForm,
    AiChatSidebar,
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
    const { alarmInfos, entityInfos, fetchInfos } = useEntityInfos(); // TODO: may be remove this

    const form = ref(metaAlarmRuleToForm(config.value.rule));
    const ruleId = computed(() => config.value.rule?._id);

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,
      ruleId,

      modal: toRef(props, 'modal'),
      context: LLM_SOCKET_CONTEXTS.metaAlarmRule,
    });

    const title = computed(() => config.value.title ?? t('modals.metaAlarmRule.create.title'));

    const { submit, submitting, submitLabel, isNew } = useSubmittableForm({
      form,
      item: config.value.rule,
      method: async () => {
        const result = await config.value.action(formToMetaAlarmRule(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useEntityInfoPropertyFetching();
    useFormConfirmableCloseModal({ form, submit, close });

    onMounted(() => {
      fetchInfos();
    });

    return {
      form,

      ruleId,
      config,
      title,

      alarmInfos,
      entityInfos,

      isNew,
      submitting,
      submitLabel,
      submit,
      close,

      chatShown,
      chatOptions,
    };
  },
};
</script>
