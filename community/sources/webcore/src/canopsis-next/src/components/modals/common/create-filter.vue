<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <patterns-form
          v-model="form"
          v-bind="patternsProps"
          autofocus
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
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
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { omit } from 'lodash';

import {
  MODALS,
  PATTERNS_FIELDS,
  SIDE_BARS,
  VALIDATION_DELAY,
  LLM_AI_CHAT_WIDTH,
} from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import PatternsForm from '@/components/forms/patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PatternsForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close, modal } = useInnerModal(props);
    const { t } = useI18n();

    /**
     * Get pattern fields based on modal configuration flags.
     *
     * Returns an array of pattern field constants that are enabled in the modal config.
     * Each flag (withAlarm, withEntity, etc.) determines which pattern field type should be included.
     *
     * @returns {string[]} Array of pattern field constants (e.g., 'alarm_pattern', 'entity_pattern')
     */
    const getPatternsFields = () => {
      const { withAlarm, withEntity, withPbehavior, withEvent, withServiceWeather } = modal.value.config;

      return [
        withAlarm && PATTERNS_FIELDS.alarm,
        withEntity && PATTERNS_FIELDS.entity,
        withPbehavior && PATTERNS_FIELDS.pbehavior,
        withEvent && PATTERNS_FIELDS.event,
        withServiceWeather && PATTERNS_FIELDS.serviceWeather,
      ].filter(Boolean);
    };

    const form = ref(filterToForm(config.value.filter, getPatternsFields()));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToFilter(form.value, getPatternsFields(), config.value.corporate));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));

    const patternsProps = computed(() => omit(config.value, ['title', 'action']));

    const sidebar = useSidebar();

    onMounted(() => sidebar.show({
      id: props.modal.id,
      name: SIDE_BARS.aiChat,
      config: {
        minimizable: true,
        width: LLM_AI_CHAT_WIDTH,
        color: 'primary',
        titleIcon: '$vuetify.icons.ai',
        titleMinimized: 'AI',
      },
    }));

    onBeforeUnmount(() => sidebar.hide({ id: props.modal.id }));

    return {
      title,
      form,
      patternsProps,
      isDisabled,
      submitting,
      close,
      submit,
    };
  },
};
</script>
