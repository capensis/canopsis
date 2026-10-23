<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ config.title ?? t('modals.applyEventFilter.title') }}</span>
      </template>

      <template #text>
        <div class="position-relative">
          <c-progress-overlay :pending="eventPatternAttributesPending" />
          <v-alert v-if="config.infoAlert" class="mb-4" type="info">
            {{ config.infoAlert }}
          </v-alert>
          <c-event-filter-patterns-field
            v-model:patterns="patterns"
            :attributes="eventPatternAttributes"
            :required="config.required ?? true"
            name="patterns"
          />
        </div>
      </template>

      <template #actions>
        <v-btn :disabled="isSubmitting" variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { formGroupsToPatternRules, patternToForm } from '@/helpers/entities/pattern/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

const promisedWait = delay => new Promise(resolve => setTimeout(resolve, delay));

export default {
  name: MODALS.applyEventFilter,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});

    const form = useForm({
      initialValues: {
        patterns: patternToForm({ event_pattern: config.value.eventPattern }),
      },
    });
    const { values, handleSubmit, isSubmitting } = form;
    const patterns = form.model('patterns');

    const { fetchEventRecordPatternFields } = usePatternsFields();

    const { pending: eventPatternAttributesPending, eventPatternAttributes } =
      usePatternsFieldsFetching(fetchEventRecordPatternFields);

    const closeModal = () => modals.hide({ id: props.modal.id });
    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(
        formGroupsToPatternRules(values.value.patterns?.groups ?? []),
      );

      await promisedWait(3000);
      await config.value.afterSubmit?.(result);
      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      config,
      patterns,
      values,
      isSubmitting,
      eventPatternAttributes,
      eventPatternAttributesPending,
      submit,
      closeModal,
    };
  },
};
</script>
