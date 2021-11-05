<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        idle-rule-form(v-model="form")
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

import { formToIdleRule, idleRuleToForm } from '@/helpers/forms/idle-rule';

import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';

import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createIdleRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    IdleRuleForm,
    ModalWrapper,
  },
  mixins: [
    validationErrorsMixin(),
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: idleRuleToForm(this.modal.config.idleRule),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createAlarmIdleRule.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToIdleRule(this.form));
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