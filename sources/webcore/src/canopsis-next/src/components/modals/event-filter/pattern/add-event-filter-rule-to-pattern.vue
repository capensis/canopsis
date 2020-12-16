<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.eventFilterRule.addAField') }}
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
import { isObject } from 'lodash';

import uid from '@/helpers/uid';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PatternRuleForm from '@/components/other/pattern/form/pattern-rule-form.vue';

import ModalWrapper from '../../modal-wrapper.vue';

function ruleToForm({ field = '', value = '' } = {}) {
  const isSimple = !isObject(value);
  const form = {
    field,
    value: '',
    advancedMode: !isSimple,
    advancedFields: [],
  };

  if (isSimple) {
    form.value = value;
  } else {
    form.advancedFields = Object.entries(value)
      .map(([fieldKey, fieldValue]) => ({ key: uid(), operator: fieldKey, value: fieldValue }));
  }

  return form;
}

function formToRule(form) {
  if (!form.advancedMode) {
    return {
      field: form.field,
      value: form.value,
    };
  }

  const value = form.advancedFields.reduce((acc, field) => {
    acc[field.operator] = field.value;

    return acc;
  }, {});

  return {
    value,

    field: form.field,
  };
}

export default {
  name: MODALS.addEventFilterRuleToPattern,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PatternRuleForm, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    const { rule = {} } = this.modal.config;

    return {
      form: ruleToForm(rule),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToRule(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

