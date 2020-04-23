<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
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
import { omit } from 'lodash';
import { MODALS } from '@/constants';

import { metaAlarmRuleToForm } from '@/helpers/forms/meta-alarm-rule';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

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
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { rule, isDuplicating } = this.modal.config;

    return {
      form: metaAlarmRuleToForm(isDuplicating ? omit(rule, ['_id']) : rule),
    };
  },
  computed: {
    title() {
      let type = 'create';

      if (this.config.rule) {
        type = this.config.isDuplicating ? 'duplicate' : 'edit';
      }

      return this.$t(`modals.metaAlarmRule.${type}.title`);
    },

    isDisabledIdField() {
      return this.config.rule && !this.config.isDuplicating;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>

