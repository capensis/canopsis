<template>
  <v-form data-vv-scope="general" @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <template-testing-test-variables-wrapper
          v-field="form"
          :is-new="isNew"
          :type="type"
        >
          <template #default="{ templateVars }">
            <event-filter-form
              v-model="form"
              :template-vars="templateVars"
              :is-disabled-id-field="config.isDisabledIdField"
            />
          </template>
        </template-testing-test-variables-wrapper>
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
import { computed, ref, watch, inject } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { eventFilterToForm, formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import {
  isChangeEntityEventFilterRuleType,
  isEnrichmentEventFilterRuleType,
} from '@/helpers/entities/event-filter/rule/entity';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/template-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EventFilterForm,
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
    const system = inject('$system');

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const { validator } = useValidationFormErrors();

    const form = ref(eventFilterToForm(config.value.rule, system.timezone));

    const isNew = computed(() => !config.value.rule?._id);
    const type = computed(() => TEMPLATE_TESTING_TEST_TYPES.eventFilter);
    const title = computed(() => config.value.title ?? t('modals.createEventFilter.create.title'));
    const isEnrichment = computed(() => isEnrichmentEventFilterRuleType(form.value.type));
    const isChangeEntity = computed(() => isChangeEntityEventFilterRuleType(form.value.type));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const isFormValid = await validator.validateAll('general');

        if (isFormValid) {
          if (config.value.action) {
            await config.value.action(formToEventFilter(form.value, system.timezone));
          }

          close();
        }
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    watch(() => form.value.type, () => validator.errors.clear());

    return {
      form,
      config,
      isNew,
      type,
      title,
      isEnrichment,
      isChangeEntity,
      isDisabled,
      submitting,
      submit,
    };
  },
};
</script>
