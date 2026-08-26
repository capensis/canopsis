<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <event-filter-form
          :rule-id="ruleId"
          :is-disabled-id-field="isDisabledIdField"
          :event-attributes="eventAttributes"
          :attributes-pending="pending"
        />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.bind" v-on="chatOptions.on" />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          variant="text"
          @click="closeModal"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitDisabled"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, inject, toRef, watch } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS } from '@/constants';

import { eventFilterToForm, formToEventFilter } from '@/helpers/entities/event-filter/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';

export default {
  name: MODALS.createEventFilter,
  components: {
    AiChatSidebar,
    EventFilterForm,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const system = inject('$system');

    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const { t } = useI18n();

    const form = useForm({
      initialValues: eventFilterToForm(config.value.rule, system.timezone),
    });
    const { values, handleSubmit, submitting, reset } = form;

    const { fetchEventFilterPatternFields } = usePatternsFields();

    const { pending, eventAttributes } = usePatternsFieldsFetching(fetchEventFilterPatternFields);

    const ruleId = computed(() => config.value.rule?._id);
    const title = computed(() => config.value.title ?? t('modals.createEventFilter.create.title'));
    const isNew = computed(() => !config.value.rule?._id || config.value.duplicating);
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const isDisabledIdField = computed(
      () => config.value.isDisabledIdField || config.value.submittingIdField,
    );

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form: values,
      modal: toRef(props, 'modal'),
      ruleId: ruleId.value,
      context: LLM_SOCKET_CONTEXTS.eventFilter,
    });
    const submitDisabled = computed(() => submitting.value || chatOptions.value.bind.pending);

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToEventFilter(values.value, system.timezone));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    watch(
      () => values.value.type,
      () => reset(values.value),
    );

    return {
      config,
      ruleId,
      title,
      pending,
      eventAttributes,
      submitting,
      submitDisabled,
      isDisabledIdField,
      chatShown,
      chatOptions,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
