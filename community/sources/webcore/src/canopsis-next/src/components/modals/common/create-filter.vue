<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <patterns-form
          v-model="form"
          v-bind="patternsProps"
          autofocus
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
import { omit } from 'lodash';
import { computed, ref, toRef } from 'vue';

import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY, LLM_SOCKET_CONTEXTS } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PatternsForm from '@/components/forms/patterns-form.vue';
import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PatternsForm, ModalWrapper, AiChatSidebar },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const patternsFields = computed(() => {
      const { withAlarm, withEntity, withPbehavior, withEvent, withServiceWeather } = config.value;

      return [
        withAlarm && PATTERNS_FIELDS.alarm,
        withEntity && PATTERNS_FIELDS.entity,
        withPbehavior && PATTERNS_FIELDS.pbehavior,
        withEvent && PATTERNS_FIELDS.event,
        withServiceWeather && PATTERNS_FIELDS.serviceWeather,
      ].filter(Boolean);
    });

    const form = ref(filterToForm(config.value.filter, patternsFields.value));

    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));
    const patternsProps = computed(() => omit(config.value, ['title', 'action']));

    const chatContext = computed(() => `${LLM_SOCKET_CONTEXTS.widgetFilter}_${config.value.widgetType}`);

    const {
      chatIds,
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: config.value?.filter?._id,
      context: chatContext,
      withoutLink: config.value?.withoutLink,
    });

    /**
     * Submits the form and calls the action callback if provided
     */
    const { submit, submitting, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        const newFilter = formToFilter(form.value, patternsFields.value, config.value.corporate);

        if (config.value?.withoutLink && chatIds.value.length) {
          newFilter.llm_chat = chatIds.value.at(-1);
        }

        const result = await config.value.action?.(newFilter);

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      patternsProps,
      submitting,
      isDisabled,
      submit,
      close,

      chatShown,
      chatOptions,
    };
  },
};
</script>
