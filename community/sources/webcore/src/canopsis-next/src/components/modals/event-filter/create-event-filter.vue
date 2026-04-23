<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <template-testing-test-variables-wrapper
          v-model="form"
          :rule-id="ruleId"
          :type="type"
        >
          <template #default="{ templateVars, copyVars }">
            <event-filter-form
              v-model="form"
              :template-vars="templateVars"
              :copy-vars="copyVars"
              :is-disabled-id-field="config.isDisabledIdField"
              :event-attributes="eventAttributes"
              :attributes-pending="pending"
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
import {
  computed,
  ref,
  watch,
  inject,
  toRef,
} from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { eventFilterToForm, formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import {
  isChangeEntityEventFilterRuleType,
  isEnrichmentEventFilterRuleType,
} from '@/helpers/entities/event-filter/rule/entity';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EventFilterForm,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.eventFilter;

    const system = inject('$system');

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const { validator } = useValidationFormErrors();

    const form = ref(eventFilterToForm(config.value.rule, system.timezone));

    const { fetchEventFilterPatternFields } = usePatternsFields();

    const {
      pending,
      eventAttributes,
    } = usePatternsFieldsFetching(fetchEventFilterPatternFields);

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.rule?._id,
      context: LLM_SOCKET_CONTEXTS.eventFilter,
    });

    const ruleId = computed(() => config.value.rule?._id);
    const title = computed(() => config.value.title ?? t('modals.createEventFilter.create.title'));
    const isEnrichment = computed(() => isEnrichmentEventFilterRuleType(form.value.type));
    const isChangeEntity = computed(() => isChangeEntityEventFilterRuleType(form.value.type));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToEventFilter(form.value, system.timezone));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    watch(() => form.value.type, () => validator.errors.clear());

    return {
      form,
      config,
      ruleId,
      type,
      title,
      pending,
      eventAttributes,
      isEnrichment,
      isChangeEntity,
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
