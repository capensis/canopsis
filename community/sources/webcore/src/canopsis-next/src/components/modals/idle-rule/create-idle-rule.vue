<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        idle-rule-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToIdleRule, idleRuleToForm } from '@/helpers/forms/idle-rule';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createIdleRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    IdleRuleForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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
        if (this.config.action) {
          await this.config.action(formToIdleRule(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
