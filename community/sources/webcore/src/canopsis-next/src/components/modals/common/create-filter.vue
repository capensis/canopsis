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
          @click="$modals.hide"
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
import { ref, computed } from 'vue';
import { omit } from 'lodash';

import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

import { useInnerModal } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';
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
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    /**
     * Gets pattern fields based on config flags
     *
     * @returns {Array} Array of pattern field constants
     */
    const getPatternsFields = () => {
      const { withAlarm, withEntity, withPbehavior, withEvent, withServiceWeather } = config.value;

      return [
        withAlarm && PATTERNS_FIELDS.alarm,
        withEntity && PATTERNS_FIELDS.entity,
        withPbehavior && PATTERNS_FIELDS.pbehavior,
        withEvent && PATTERNS_FIELDS.event,
        withServiceWeather && PATTERNS_FIELDS.serviceWeather,
      ].filter(Boolean);
    };

    const form = ref(filterToForm(config.value.filter, getPatternsFields()));

    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));
    const patternsProps = computed(() => omit(config.value, ['title', 'action']));

    /**
     * Submits the form and calls the action callback if provided
     */
    const { submit, submitting, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToFilter(form.value, getPatternsFields(), config.value.corporate));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close,
    });

    return {
      form,
      title,
      patternsProps,
      submitting,
      isDisabled,
      submit,
    };
  },
};
</script>
