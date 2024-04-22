<template>
  <v-layout column>
    <c-alert
      v-show="!actions.length"
      type="info"
    >
      {{ $t('scenario.emptyActions') }}
    </c-alert>
    <c-card-iterator-field
      v-field="actions"
      :draggable-group="draggableGroup"
      class="mb-2"
      item-key="key"
    >
      <template #item="{ item: action, index }">
        <scenario-action-field
          v-field="actions[index]"
          :name="`${name}.${action.key}`"
          :action-number="index + 1"
          :has-previous-webhook="hasPreviousWebhook(index)"
          @remove="removeItemFromArray(index)"
          @duplicate="duplicateAction(action)"
        />
      </template>
    </c-card-iterator-field>
    <c-btn-with-error
      :error="hasActionsErrors ? $t('scenario.errors.actionRequired') : ''"
      outlined
      @click="addAction"
    >
      {{ $t('scenario.addAction') }}
    </c-btn-with-error>
  </v-layout>
</template>

<script>
import { cloneDeep } from 'lodash';

import { actionToForm, isWebhookActionType } from '@/helpers/entities/action';
import { uid } from '@/helpers/uid';

import { formArrayMixin, validationChildrenMixin } from '@/mixins/form';

import ScenarioActionField from './fields/scenario-action-field.vue';

export default {
  inject: ['$validator'],
  components: { ScenarioActionField },
  mixins: [formArrayMixin, validationChildrenMixin],
  model: {
    prop: 'actions',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    actions: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'actions',
    },
  },
  computed: {
    hasActionsErrors() {
      return this.errors.has(this.name);
    },

    draggableGroup() {
      return {
        name: 'scenarios-actions',
      };
    },

    webhookIndexes() {
      return this.actions.reduce((acc, action, index) => {
        if (isWebhookActionType(action.type)) {
          acc.push(index);
        }

        return acc;
      }, []);
    },
  },
  watch: {
    actions() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.attachMinValueRule();
  },
  beforeDestroy() {
    this.detachMinValueRule();
  },
  methods: {
    attachMinValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'min_value:1',
        getter: () => this.actions.length,
        vm: this,
      });
    },

    detachMinValueRule() {
      this.$validator.detach(this.name);
    },

    addAction() {
      this.addItemIntoArray(actionToForm());
    },

    hasPreviousWebhook(index) {
      return this.webhookIndexes.indexOf(index) > 0;
    },

    duplicateAction(action) {
      const clonedAction = cloneDeep(action);

      this.addItemIntoArray({ ...clonedAction, key: uid() });
    },
  },
};
</script>
