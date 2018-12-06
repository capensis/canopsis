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
      h2 Actions
      v-list
        draggable(v-model="actions", @start="drag=true", @end="drag=false")
          v-list-group.grey.white--text(v-for="(action, index) in actions", :key="action.name")
            v-list-tile(slot="activator")
              v-list-tile-title {{index + 1}} - {{ action.type }} - {{ action.name || action.from }}
            v-list-tile
              v-layout(column)
                div(v-if="action.name") Name: {{ action.name }}
                div(v-if="action.value") Value: {{ action.value }}
                div(v-if="action.description") Description: {{ action.description }}
                div(v-if="action.from") From: {{ action.from }}
                div(v-if="action.to") To: {{ action.to }}
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import Draggable from 'vuedraggable';
import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.eventFilterRuleActions,
  components: {
    Draggable,
  },
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
  mounted() {
    this.actions = [...this.config.actions];
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
