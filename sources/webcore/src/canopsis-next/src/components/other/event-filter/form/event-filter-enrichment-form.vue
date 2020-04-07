<template lang="pug">
  div
    v-container
      v-divider
      h3.my-2 {{ $t('modals.eventFilterRule.enrichmentOptions') }}
      v-btn(@click="showEditActionsModal") {{ $t('modals.eventFilterRule.editActions') }}
      v-btn(@click="showEditExternalDataModal") {{ $t('modals.eventFilterRule.externalData') }}
      v-select(
        v-field="form.onSuccess",
        :label="$t('modals.eventFilterRule.onSuccess')",
        :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES)"
      )
      v-select(
        v-field="form.onFailure",
        :label="$t('modals.eventFilterRule.onFailure')",
        :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES)"
      )
    v-alert(:value="errors.has('actions')", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';

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
  },
  async created() {
    this.$validator.attach({
      name: 'actions',
      rules: 'required:true',
      getter: () => this.form.actions,
      context: () => this,
    });
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
          value: this.form.externalData,
          action: value => this.updateField('externalData', value),
        },
      });
    },
  },
};
</script>

