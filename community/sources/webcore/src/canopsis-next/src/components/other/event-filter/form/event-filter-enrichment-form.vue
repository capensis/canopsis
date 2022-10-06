<template lang="pug">
  div
    v-layout
      v-btn.mx-0(@click="showEditActionsModal") {{ $t('eventFilter.editActions') }}
    v-select(
      v-field="form.on_success",
      :label="$t('eventFilter.onSuccess')",
      :items="eventFilterAfterTypes"
    )
    v-select(
      v-field="form.on_failure",
      :label="$t('eventFilter.onFailure')",
      :items="eventFilterAfterTypes"
    )
    v-alert(:value="errors.has(name)", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES, MODALS } from '@/constants';

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
    eventFilterAfterTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES);
    },
  },
  created() {
    this.attachActionsRequiredRule();
  },
  beforeDestroy() {
    this.detachActionsRequiredRule();
  },
  methods: {
    showEditActionsModal() {
      this.$modals.show({
        name: MODALS.eventFilterActions,
        config: {
          actions: this.form.actions,
          action: (updatedActions) => {
            this.updateField('actions', updatedActions);
            this.$nextTick(() => this.$validator.validate('actions'));
          },
        },
      });
    },

    attachActionsRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => this.form.actions,
      });
    },

    detachActionsRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
