<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <declare-ticket-rule-form
          v-model="form"
          :rule-id="config.declareTicketRule?._id"
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
import { computed, ref } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { declareTicketRuleToForm, formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDeclareTicketRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DeclareTicketRuleForm,
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

    const form = ref(declareTicketRuleToForm(config.value.declareTicketRule));

    const isNew = computed(() => !config.value.declareTicketRule?._id);
    const title = computed(() => config.value.title ?? t('modals.createDeclareTicketRule.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const declareTicketRule = await config.value.action?.(formToDeclareTicketRule(form.value));

        close();

        return declareTicketRule;
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
