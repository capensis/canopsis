<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.eventFilterRule.addAction') }}
    v-card-text
      v-form(ref="form")
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
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP } from '@/constants';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRightMixin from '@/mixins/entities/right';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
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
    computed: {
      getComponentByOption() {
        return (option = {}) => (option.value === 'value' ? 'mixed-field' : 'v-text-field');
      },
      getValidationRulesByOption() {
        return (option = {}) => option.required && 'required';
      },
    },
  },
  mounted() {
    if (this.config.ruleAction) {
      const { ruleAction } = this.config;

      this.form = {
        ...this.form,

        type: EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP[ruleAction.type],
      };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          type: this.form.type.value,
          ...pick(this.form, Object.keys(this.form.type.options)),
        };

        if (this.config.action) {
          await this.config.action(data);
        }

        this.addSuccessPopup({ text: this.$t('success.default') });
        this.hideModal();
      }
    },
  },
};
</script>
