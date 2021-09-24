<template lang="pug">
  v-layout(column)
    v-layout(align-center, justify-space-between)
      h2 {{ $t('modals.eventFilterRule.actions') }}
      c-action-btn(
        icon="add",
        color="primary",
        :tooltip="$t('common.add')",
        @click="showCreateActionModal"
      )
    v-list.py-0(dark)
      draggable(v-field="actions", :class="{ 'pa-2': !actions.length }")
        event-filter-rule-actions-list-item(
          v-for="(action, index) in actions",
          :key="`${action.type}_${action.name}`",
          :action="action",
          :action-number="index + 1",
          @edit="showEditActionModal(index, $event)",
          @remove="removeItemFromArray(index)"
        )
</template>

<script>
import Draggable from 'vuedraggable';
import { MODALS } from '@/constants';

import { formArrayMixin } from '@/mixins/form/array';

import EventFilterRuleActionsListItem from './partials/event-filter-rule-actions-list-item.vue';

export default {
  components: { Draggable, EventFilterRuleActionsListItem },
  mixins: [
    formArrayMixin,
  ],
  model: {
    prop: 'actions',
    event: 'input',
  },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    showCreateActionModal() {
      this.$modals.show({
        name: MODALS.createEventFilterRuleAction,
        config: {
          action: action => this.addItemIntoArray(action),
        },
      });
    },

    showEditActionModal(index, ruleAction) {
      this.$modals.show({
        name: MODALS.createEventFilterRuleAction,
        config: {
          ruleAction,
          title: this.$t('modals.eventFilterRule.editAction'),
          action: action => this.updateItemInArray(index, action),
        },
      });
    },
  },
};
</script>
