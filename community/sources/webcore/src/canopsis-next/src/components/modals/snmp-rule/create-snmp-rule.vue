<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <snmp-rule-form v-model="form" />
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
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { snmpRuleToForm, formToSnmpRule } from '@/helpers/entities/snmp-rule/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createSnmpRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    SnmpRuleForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const {
      config,
      close,
    } = useInnerModal(props);

    const form = ref(snmpRuleToForm(config.value.snmpRule));

    const {
      submit,
      isDisabled,
      submitting,
    } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToSnmpRule(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close,
    });

    const title = computed(() => config.value.title || t('modals.createSnmpRule.create.title'));

    return {
      form,

      isDisabled,
      submitting,

      title,

      submit,
      close,
    };
  },
};
</script>
