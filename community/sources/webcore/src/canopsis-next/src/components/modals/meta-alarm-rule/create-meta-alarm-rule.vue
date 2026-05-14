<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <v-layout class="gap-3" column>
          <c-enabled-field v-model="form.enabled" hide-details with-background />
          <c-form-general-patterns-tabs
            v-model="form"
            :rule-id="ruleId"
            :type="type"
            :patterns-label="$t('metaAlarmRule.patternsTabLabel')"
          >
            <template #general="{ setRef, templateVars }">
              <meta-alarm-rule-general-form
                v-model="form"
                :ref="setRef"
                :disabled-id-field="config.isDisabledIdField"
                :template-vars="templateVars"
              />
            </template>
            <template #patterns="{ setRef, templateVars }">
              <meta-alarm-rule-parameters-form
                v-model="form"
                :ref="setRef"
                :template-vars="templateVars"
              />
            </template>
          </c-form-general-patterns-tabs>
        </v-layout>
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
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, onMounted, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfos } from '@/hooks/store/modules/entity-infos';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import MetaAlarmRuleGeneralForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue';
import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    MetaAlarmRuleGeneralForm,
    MetaAlarmRuleParametersForm,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule;

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

    const { submit, isDisabled, submitting, submitLabel } = useSubmittableForm({
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
      type,

      form,

      ruleId,
      config,
      title,

      alarmInfos,
      entityInfos,

      isDisabled,
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
