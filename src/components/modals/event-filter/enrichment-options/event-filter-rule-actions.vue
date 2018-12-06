<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Event filter actions
    v-card-text
      v-card
        v-card-title.primary.white--text Add action
        v-card-text
          v-select(
          :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES)",
          v-model="actionForm.type",
          return-object,
          item-text="value",
          label="Type",
          )
          v-text-field(
          v-for="option in actionForm.type.options",
          :key="option.value"
          v-model="actionForm[option.value]",
          :label="option.text",
          hide-details,
          :required="isRequired(actionForm.type, option)"
          )
        v-divider
        v-btn.primary(@click="addAction") Add
    v-container
      div(v-for="action in actions") {{ action }}
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.eventFilterRuleActions,
  mixins: [modalInnerMixin],
  data() {
    return {
      actions: [],
      actionForm: {
        type: EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
        name: '',
        value: '',
        description: '',
        from: '',
        to: '',
      },
    };
  },
  methods: {
    isRequired(actionType, option) {
      return actionType.options[option.value].required;
    },
    addAction() {
      const action = {
        type: this.actionForm.type.value,
      };

      Object.keys(this.actionForm.type.options).forEach(option => action[option] = this.actionForm[option]);

      this.actions.push(action);
    },
    submit() {
      this.config.action(this.actions);
    },
  },
};
</script>
