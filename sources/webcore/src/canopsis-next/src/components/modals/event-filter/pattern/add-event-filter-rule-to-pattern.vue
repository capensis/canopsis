<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.eventFilterRule.addAField') }}
    template(slot="text")
      add-event-filter-rule-to-pattern-form(
        v-model="form",
        :operators="config.operators",
        :isSimple="config.isSimple"
      )
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { isObject } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import AddEventFilterRuleToPatternForm
  from '@/components/other/event-filter/form/add-event-filter-rule-to-pattern-form.vue';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.addEventFilterRuleToPattern,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AddEventFilterRuleToPatternForm, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        advancedMode: false,
        field: '',
        value: '',
        advancedRuleFields: [],
      },
    };
  },
  mounted() {
    if (this.config) {
      const {
        operators,
        ruleKey = '',
        ruleValue = '',
      } = this.config;

      const isSimpleRule = !isObject(ruleValue);

      this.operators = operators;
      this.form.advancedMode = !isSimpleRule;
      this.form.field = ruleKey;

      if (isSimpleRule) {
        this.form.value = ruleValue;
      } else {
        this.form.advancedRuleFields = Object.entries(ruleValue).map(([key, value]) => ({ key, value }));
      }
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          let newRule = {};

          if (!this.form.advancedMode) {
            newRule = { field: this.form.field, value: this.form.value };
          } else {
            const value = this.form.advancedRuleFields.reduce((acc, field) => {
              acc[field.key] = field.value;
              return acc;
            }, {});
            newRule = { field: this.form.field, value };
          }

          await this.config.action(newRule);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

