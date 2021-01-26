<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        meta-alarm-rule-form(
          v-model="form",
          :isDisabledIdField="isDisabledIdField"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/forms/meta-alarm-rule';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    MetaAlarmRuleForm,
    ModalWrapper,
  },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { rule } = this.modal.config;

    return {
      form: metaAlarmRuleToForm(rule),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.metaAlarmRule.create.title');
    },

    isDisabledIdField() {
      return this.config.isDisabledIdField;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToMetaAlarmRule(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>

