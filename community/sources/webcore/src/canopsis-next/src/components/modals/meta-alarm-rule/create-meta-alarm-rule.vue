<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        meta-alarm-rule-form(
          v-model="form",
          :is-disabled-id-field="isDisabledIdField"
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: metaAlarmRuleToForm(this.modal.config.rule),
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
