<template lang="pug">
   v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      v-form
        v-select(:items="ruleTypes", v-model="form.type", :label="$t('common.type')")
        v-text-field(v-model.number="form.priority", type="number", :label="$t('modals.eventFilterRule.priority')")
        v-switch(v-model="form.enabled", :label="$t('common.enabled')")
      v-btn(@click="editPattern") {{ $t('modals.eventFilterRule.editPattern') }}
      template(v-if="form.type === $constants.EVENT_FILTER_RULE_TYPES.enrichment")
        v-container
          v-divider
          h3.my-2 {{ $t('modals.eventFilterRule.enrichmentOptions') }}
          v-btn(@click="editActions") {{ $t('modals.eventFilterRule.editActions') }}
          v-btn(@click="editExternalData") {{ $t('modals.eventFilterRule.externalData') }}
          v-select(
          :label="$t('modals.eventFilterRule.onSuccess')",
          v-model="enrichmentOptions.onSuccess",
          :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES)",
          )
          v-select(
          :label="$t('modals.eventFilterRule.onFailure')",
          v-model="enrichmentOptions.onFailure",
          :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES)",
          )
    v-divider
    v-alert(:value="errors.has('actions')", type="error") {{ $t('eventFilter.actionsRequired') }}
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import { MODALS, EVENT_FILTER_RULE_TYPES, EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.createEventFilterRule,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      ruleTypes: Object.values(EVENT_FILTER_RULE_TYPES),
      form: {
        type: EVENT_FILTER_RULE_TYPES.drop,
        pattern: {},
        priority: 0,
        enabled: true,
      },
      enrichmentOptions: {
        actions: [],
        externalData: {},
        onSuccess: EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
        onFailure: EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
      },
    };
  },
  mounted() {
    if (this.config.rule) {
      const {
        type,
        pattern,
        priority,
        enabled,
        actions,
        externalData,
        on_success: onSuccess,
        on_failure: onFailure,
      } = cloneDeep(this.config.rule);

      this.form = {
        type,
        pattern,
        priority,
        enabled,
      };

      this.enrichmentOptions = {
        actions,
        externalData,
        onSuccess: onSuccess || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
        onFailure: onFailure || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
      };
    }
  },
  created() {
    this.$validator.attach('actions', 'required', {
      getter: () => this.enrichmentOptions.actions,
      context: () => this,
    });
  },
  methods: {
    editPattern() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          pattern: this.form.pattern,
          action: pattern => this.form.pattern = pattern,
        },
      });
    },
    editActions() {
      this.showModal({
        name: MODALS.eventFilterRuleActions,
        config: {
          actions: this.enrichmentOptions.actions,
          action: updatedActions => this.enrichmentOptions.actions = updatedActions,
        },
      });
    },
    editExternalData() {
      this.showModal({
        name: MODALS.eventFilterRuleExternalData,
        config: {
          value: this.enrichmentOptions.externalData,
          action: value => this.enrichmentOptions.externalData = value,
        },
      });
    },
    async submit() {
      if (this.form.type === 'enrichment') {
        const isFormValid = await this.$validator.validate('actions');
        if (isFormValid) {
          this.config.action({
            ...this.form,
            actions: this.enrichmentOptions.actions,
            external_data: this.enrichmentOptions.externalData,
            on_success: this.enrichmentOptions.onSuccess,
            on_failure: this.enrichmentOptions.onFailure,
          });
          this.hideModal();
        }
      } else {
        this.config.action({ ...this.form });
        this.hideModal();
      }
    },
  },
};
</script>

