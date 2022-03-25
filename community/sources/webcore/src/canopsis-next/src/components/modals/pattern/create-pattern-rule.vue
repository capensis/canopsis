<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('eventFilter.addAField') }}
      template(slot="text")
        pattern-rule-form(
          v-model="form",
          :operators="config.operators",
          :only-simple="config.onlySimple"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { patternRuleToForm } from '@/helpers/forms/pattern';

import { modalInnerMixin } from '@/mixins/modal/inner';

import PatternRuleForm from '@/components/other/pattern/form/pattern-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPatternRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PatternRuleForm, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    const { rule = {} } = this.modal.config;

    return {
      form: patternRuleToForm(rule),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
