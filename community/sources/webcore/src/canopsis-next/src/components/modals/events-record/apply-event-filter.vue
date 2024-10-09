<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title ?? $t('modals.applyEventFilter.title') }}</span>
      </template>
      <template #text="">
        <c-event-filter-patterns-field
          v-model="form"
          :excluded-attributes="config.excludedAttributes"
          name="patterns"
          required
          @input="errors.remove('patterns')"
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
import { ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formGroupsToPatternRules, patternToForm } from '@/helpers/entities/pattern/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.applyEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(patternToForm({ event_pattern: config.value.eventPattern }));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formGroupsToPatternRules(form.value?.groups));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      config,
      form,
      isDisabled,
      submitting,

      submit,
      close,
    };
  },
};
</script>
