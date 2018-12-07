<template lang="pug">
   v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Event filter rule
    v-card-text
      v-form
        v-select(:items="ruleTypes", v-model="form.type", label="Type")
        v-text-field(v-model="form.priority", type="number", label="Priority")
        v-switch(v-model="form.enabled", label="Enabled")
      v-btn(@click="editPattern") Edit pattern
      template(v-if="form.type === this.$constants.EVENT_FILTER_RULE_TYPES.enrichment")
        v-container
          v-divider
          h3.my-2 Enrichment options
          v-btn(@click="editActions") Edit actions
          v-btn(@click="editExternalData") Exernal data
          v-select(label="On success", v-model="enrichmentOptions.on_success", :items=['pass', 'break', 'drop'])
          v-select(label="On failure", v-model="enrichmentOptions.on_failure", :items=['pass', 'break', 'drop'])
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, EVENT_FILTER_RULE_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.createEventFilterRule,
  mixins: [modalInnerMixin],
  data() {
    return {
      ruleTypes: Object.values(EVENT_FILTER_RULE_TYPES),
      form: {
        type: EVENT_FILTER_RULE_TYPES.drop,
        pattern: {
          component: 'canopsis',
          state: { '<=': 3, '>=': 1 },
        },
        priority: 0,
        enabled: true,
      },
      enrichmentOptions: {
        actions: [],
        external_data: {},
        on_success: 'pass',
        on_failure: 'pass',
      },
    };
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
          value: this.enrichmentOptions.external_data,
          action: value => this.enrichmentOptions.external_data = value,
        },
      });
    },
    submit() {
      if (this.form.type === 'enrichment') {
        this.config.action({ ...this.form, ...this.enrichmentOptions });
      } else {
        this.config.action({ ...this.form });
      }

      this.hideModal();
    },
  },
};
</script>

