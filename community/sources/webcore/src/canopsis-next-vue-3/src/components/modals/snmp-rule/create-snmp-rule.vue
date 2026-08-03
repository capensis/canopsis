<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <snmp-rule-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
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
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { formToSnmpRule, snmpRuleToForm } from '@/helpers/entities/snmp-rule/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

export default {
  name: MODALS.createSnmpRule,

  components: {
    SnmpRuleForm,
  },

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
    const close = () => modals.hide({ id: props.modal.id });

    const form = useForm({
      initialValues: snmpRuleToForm(config.value.snmpRule),
    });

    const title = computed(() => config.value.title || t('modals.createSnmpRule.create.title'));

    const { submitting, submit, submitLabel } = useSubmittableForm({
      item: config.value.snmpRule,
      method: form.handleSubmit(async () => {
        await config.value.action?.(formToSnmpRule(form.values.value));

        close();
      }),
    });

    useFormConfirmableCloseModal({
      form: form.values,
      submit,
      close,
    });

    return {
      t,
      title,
      submitting,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>
