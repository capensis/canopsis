<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        alarm-status-rule-form(v-model="form", :flapping="config.flapping")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { alarmStatusRuleToForm, formToAlarmStatusRule } from '@/helpers/forms/alarm-status-rule';

import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';
import { validationErrorsMixin } from '@/mixins/form/validation-errors';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAlarmStatusRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmStatusRuleForm, ModalWrapper },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
    validationErrorsMixin(),
  ],
  data() {
    const { rule, flapping } = this.modal.config;

    return {
      form: alarmStatusRuleToForm(rule, flapping),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToAlarmStatusRule(this.form));
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
