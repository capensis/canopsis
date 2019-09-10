<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.eventFilterRule.addAction') }}
    v-card-text
      v-form
        v-select(
          v-model="form.type",
          :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES)",
          :label="$t('common.type')",
          item-text="value",
          return-object
        )
        component(
          v-for="option in form.type.options"
          v-model="form[option.value]",
          v-validate="getValidationRulesByOption(option)",
          :is="getComponentByOption(option)",
          :key="option.value",
          :label="option.text",
          :name="option.value",
          :error-messages="errors.collect(option.value)"
        )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP } from '@/constants';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRightMixin from '@/mixins/entities/right';

import MixedField from '@/components/forms/fields/mixed-field.vue';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
  filters: {
    ruleActionToForm(ruleAction) {
      const type = EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP[ruleAction.type];

      return {
        ...ruleAction,

        type,
      };
    },
    formToRuleAction(form) {
      return {
        ...pick(form, Object.keys(form.type.options)),

        type: form.type.value,
      };
    },
  },
  components: { MixedField },
  mixins: [popupMixin, modalInnerMixin, entitiesRightMixin],
  data() {
    const enrichmentActionsTypes = cloneDeep(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES);

    return {
      actions: [],
      editableActionId: null,
      form: {
        type: enrichmentActionsTypes.setField,
        name: '',
        value: '',
        description: '',
        from: '',
        to: '',
      },
    };
  },
  computed: {
    getComponentByOption() {
      return (option = {}) => (option.value === 'value' ? 'mixed-field' : 'v-text-field');
    },
    getValidationRulesByOption() {
      return (option = {}) => option.required && 'required';
    },
  },
  mounted() {
    if (this.config.ruleAction) {
      const { ruleAction } = this.config;

      this.form = {
        ...this.form,
        ...this.$options.filters.ruleActionToForm(ruleAction),
      };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = this.$options.filters.formToRuleAction(this.form);

        if (this.config.action) {
          await this.config.action(data);
        }

        this.hideModal();
      }
    },
  },
};
</script>
