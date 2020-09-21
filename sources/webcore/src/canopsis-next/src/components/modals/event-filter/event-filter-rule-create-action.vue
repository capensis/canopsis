<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.eventFilterRule.addAction') }}
      template(slot="text")
        div
          v-select(
            v-model="form.type",
            :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES)",
            :label="$t('common.type')",
            item-text="value",
            return-object
          )
          component(
            v-for="option in form.type.options",
            v-model="form[option.value]",
            v-validate="$options.filters.validationRulesByOption(option)",
            :is="option | componentByOption",
            :key="option.value",
            :label="option.text",
            :name="option.value",
            :error-messages="errors.collect(option.value)"
          )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesRightMixin from '@/mixins/entities/right';

import MixedField from '@/components/forms/fields/mixed-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventFilterRuleCreateAction,
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
    componentByOption({ value } = {}) {
      return value === 'value' ? 'mixed-field' : 'v-text-field';
    },
    validationRulesByOption({ required } = {}) {
      return required && 'required';
    },
  },
  components: { MixedField, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin(), entitiesRightMixin],
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

        this.$modals.hide();
      }
    },
  },
};
</script>
