<template>
  <c-card-iterator-form
    v-field="actions"
    :draggable-group="draggableGroup"
    :name="name"
    :required-error-message="$t('scenario.errors.actionRequired')"
    :empty-message="$t('scenario.emptyActions')"
    :add-button-label="$t('scenario.addAction')"
    item-key="key"
    iterator-class="mb-2"
    required
    @add="addAction"
  >
    <template #item="{ item: action, index }">
      <scenario-action-field
        v-field="actions[index]"
        :name="`${name}.${action.key}`"
        :action-number="index + 1"
        :has-previous-webhook="hasPreviousWebhook(index)"
        :template-vars="templateVars"
        @remove="removeItemFromArray(index)"
        @duplicate="duplicateAction(action)"
      />
    </template>
  </c-card-iterator-form>
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
    templateVars: {
      type: Object,
      default: () => ({}),
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

    firstWebhookIndex() {
      return this.webhookIndexes[0];
    },
  },
  methods: {
    addAction() {
      this.addItemIntoArray(actionToForm());
    },

    hasPreviousWebhook(index) {
      return this.firstWebhookIndex < index;
    },

    duplicateAction(action) {
      const clonedAction = cloneDeep(action);

      this.addItemIntoArray({ ...clonedAction, key: uid() });
    },
  },
};
</script>
