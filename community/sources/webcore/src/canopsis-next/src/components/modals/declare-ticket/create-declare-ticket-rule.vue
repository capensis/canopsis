<template>
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
          <v-form>
            <declare-ticket-rule-form
              v-model="form"
              :template-vars="templateVars"
            />
          </v-form>
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
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import {
  declareTicketRuleToForm,
  formToDeclareTicketRule,
  declareTicketRuleErrorsToForm,
} from '@/helpers/entities/declare-ticket/rule/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';
import { usePopups } from '@/hooks/popups';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDeclareTicketRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DeclareTicketRuleForm,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const { setFormErrors } = useValidationFormErrors();
    const popups = usePopups();

    const form = ref(declareTicketRuleToForm(config.value.declareTicketRule));

    const isNew = computed(() => !config.value.declareTicketRule?._id);
    const title = computed(() => config.value.title ?? t('modals.createDeclareTicketRule.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        try {
          await config.value.action?.(formToDeclareTicketRule(form.value));

          close();
        } catch (err) {
          if (err.error) {
            popups.error({ text: err.error });
          } else {
            setFormErrors(declareTicketRuleErrorsToForm(err, form.value));
          }
        }
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      isDisabled,
      submitting,
      submit,
    };
  },
};
</script>
