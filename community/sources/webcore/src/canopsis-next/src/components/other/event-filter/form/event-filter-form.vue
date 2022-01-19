<template lang="pug">
  div
    c-id-field(v-field="form._id", :disabled="isDisabledIdField", :help-text="$t('eventFilter.idHelp')")
    c-event-filter-type-field(v-field="form.type")
    c-description-field(v-field="form.description", required)
    c-priority-field(v-field="form.priority")
    c-enabled-field(v-field="form.enabled")
    patterns-list(v-field="form.patterns")

    template(v-if="isChangeEntityType")

    template(v-if="isEnrichmentType")
      v-container.pa-0
        v-divider
        h3.my-2 {{ $t('eventFilter.enrichmentOptions') }}
        v-layout
          v-btn.mx-0(@click="showEditActionsModal") {{ $t('eventFilter.editActions') }}
          v-btn(@click="showEditExternalDataModal") {{ $t('eventFilter.externalData') }}
        v-select(
          v-field="form.config.on_success",
          :label="$t('eventFilter.onSuccess')",
          :items="eventFilterAfterTypes"
        )
        v-select(
          v-field="form.config.on_failure",
          :label="$t('eventFilter.onFailure')",
          :items="eventFilterAfterTypes"
        )
      v-alert(:value="errors.has('actions')", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES, EVENT_FILTER_RULE_TYPES, MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
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
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isEnrichmentType() {
      return this.form.type === EVENT_FILTER_RULE_TYPES.enrichment;
    },

    isChangeEntityType() {
      return this.form.type === EVENT_FILTER_RULE_TYPES.changeEntity;
    },

    eventFilterAfterTypes() {
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
          actions: this.form.config.actions,
          action: (actions) => {
            this.updateField('config.actions', actions);
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
          action: externalData => this.updateField('external_data', externalData),
        },
      });
    },

    attachRequiredRule() {
      this.$validator.attach({
        name: 'actions',
        rules: 'required:true',
        getter: () => (this.isEnrichmentType ? this.form.actions : true),
        context: () => this,
      });
    },

    detachRequiredRule() {
      this.$validator.detach('actions');
    },
  },
};
</script>
