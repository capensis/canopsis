<template lang="pug">
  div
    v-container.pa-0
      v-divider
      h3.my-2 {{ $t('modals.eventFilterRule.enrichmentOptions') }}
      v-layout
        v-btn.mx-0(@click="showEditActionsModal") {{ $t('modals.eventFilterRule.editActions') }}
        v-btn(@click="showEditExternalDataModal") {{ $t('modals.eventFilterRule.externalData') }}
      v-select(
        v-field="form.on_success",
        :label="$t('modals.eventFilterRule.onSuccess')",
        :items="successItems"
      )
      v-select(
        v-field="form.on_failure",
        :label="$t('modals.eventFilterRule.onFailure')",
        :items="failureItems"
      )
    v-alert(:value="errors.has('actions')", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES, MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'actions',
    },
  },
  computed: {
    failureItems() {
      return Object.values(EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES);
    },

    successItems() {
      return Object.values(EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES);
    },
  },
  created() {
    this.attachRequiredRule();
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    showEditActionsModal() {
      this.$modals.show({
        name: MODALS.eventFilterRuleActions,
        config: {
          actions: this.form.actions,
          action: (updatedActions) => {
            this.updateField('actions', updatedActions);
            this.$nextTick(() => this.$validator.validate('actions'));
          },
        },
      });
    },

    showEditExternalDataModal() {
      this.$modals.show({
        name: MODALS.eventFilterRuleExternalData,
        config: {
          value: this.form.external_data,
          action: value => this.updateField('external_data', value),
        },
      });
    },

    attachRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => this.form.actions,
        context: () => this,
      });
    },

    detachRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
