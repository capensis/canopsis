<template>
  <c-card-iterator-form
    :value="actionsField"
    :handle="`.${dragItemHandleClass}`"
    :required-error-message="t('scenario.errors.actionRequired')"
    :empty-message="t('scenario.emptyActions')"
    :add-button-label="t('scenario.addAction')"
    item-key="key"
    name="actions"
    required
    @update:value="actionsField = $event"
    @add="addAction"
  >
    <template #item="{ index, item: action }">
      <scenario-action-field
        :key="action.key"
        :model-value="actionsField[index]"
        :name="`actions.${action.key}`"
        :action-number="index + 1"
        :has-previous-webhook="hasPreviousWebhook(index)"
        :template-vars="templateVars"
        @update:model-value="updateAction(index, $event)"
        @remove="removeAction(index)"
        @duplicate="duplicateAction(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { cloneDeep } from 'lodash';

import { actionToForm, isWebhookActionType } from '@/helpers/entities/action/form';
import { uid } from '@/helpers/uid';

import { useFormArray } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ScenarioActionField from './fields/scenario-action-field.vue';

export default {
  components: {
    ScenarioActionField,
  },

  props: {
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const actionsField = useFormArray('actions');
    const dragItemHandleClass = 'scenario-action-drag-handle';

    const addAction = () => {
      actionsField.value = [...actionsField.value, actionToForm()];
    };

    const updateAction = (index, action) => {
      const actions = [...actionsField.value];

      actions[index] = action;
      actionsField.value = actions;
    };

    const removeAction = (index) => {
      actionsField.value = actionsField.value.filter((action, actionIndex) => actionIndex !== index);
    };

    const duplicateAction = (index) => {
      const action = cloneDeep(actionsField.value[index]);

      action.key = uid();
      actionsField.value = [...actionsField.value, action];
    };

    const hasPreviousWebhook = index =>
      actionsField.value
        .slice(0, index)
        .some(action => isWebhookActionType(action.type));

    return {
      t,
      actionsField,
      dragItemHandleClass,
      addAction,
      updateAction,
      removeAction,
      duplicateAction,
      hasPreviousWebhook,
    };
  },
};
</script>
