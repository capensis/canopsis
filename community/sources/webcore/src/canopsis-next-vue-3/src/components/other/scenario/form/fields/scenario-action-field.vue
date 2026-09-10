<template>
  <c-card-iterator-item
    :item-number="actionNumber"
    :drag-handle-class="dragHandleClass"
    :removable="false"
  >
    <template #header>
      <c-layout gap="2" align="start">
        <c-action-type-field
          :model-value="modelValue.type"
          :name="`${name}.type`"
          required
          @update:model-value="updateField('type', $event)"
        />

        <c-layout class="scenario-action-field__actions" gap="1" align="center">
          <c-action-btn type="duplicate" @click="duplicate" />
          <c-action-btn type="delete" @click="remove" />
        </c-layout>
      </c-layout>
    </template>

    <c-layout gap="3" column>
      <v-row density="compact">
        <v-col cols="6">
          <c-enabled-field
            :model-value="modelValue.emit_trigger_success"
            :label="t('common.emitTriggerSuccess')"
            :name="`${name}.emit_trigger_success`"
            @update:model-value="updateField('emit_trigger_success', $event)"
          />
        </v-col>

        <v-col v-if="isWebhookAction && parameters.declare_ticket?.enabled" cols="6">
          <c-enabled-field
            :model-value="modelValue.emit_trigger_fail"
            :label="t('common.emitTriggerFail')"
            :name="`${name}.emit_trigger_fail`"
            @update:model-value="updateField('emit_trigger_fail', $event)"
          />
        </v-col>
      </v-row>

      <action-author-field
        :model-value="parameters"
        :name="parametersName"
        :variables="templateVars.author"
        @update:model-value="parameters = $event"
      />

      <v-expand-transition>
        <c-form-block v-if="isWebhookAction">
          <c-form-block-row :label="t('scenario.skip')">
            <c-layout column>
              <c-enabled-field
                :model-value="parameters.skip_for_child"
                :label="t('scenario.skipForChild')"
                :name="`${parametersName}.skip_for_child`"
                @update:model-value="updateParametersField('skip_for_child', $event)"
              />

              <c-enabled-field
                :model-value="parameters.skip_for_instruction"
                :label="t('scenario.skipForInstruction')"
                :name="`${parametersName}.skip_for_instruction`"
                @update:model-value="updateParametersField('skip_for_instruction', $event)"
              />
            </c-layout>
          </c-form-block-row>
        </c-form-block>
      </v-expand-transition>

      <c-form-block>
        <c-form-block-row :label="t('scenario.workflow')">
          <c-workflow-field
            :model-value="modelValue.drop_scenario_if_not_matched"
            :stop-label="t('common.stop')"
            :continue-label="t('scenario.remainingAction')"
            :name="`${name}.drop_scenario_if_not_matched`"
            @update:model-value="updateField('drop_scenario_if_not_matched', $event)"
          />
        </c-form-block-row>

        <template v-if="isWebhookAction">
          <c-form-block-row :label="t('scenario.workflowInCaseOfFailure')">
            <c-workflow-field
              :model-value="parameters.stop_on_fail"
              :stop-label="t('common.stop')"
              :continue-label="t('scenario.remainingStep')"
              :name="`${parametersName}.stop_on_fail`"
              @update:model-value="updateParametersField('stop_on_fail', $event)"
            />
          </c-form-block-row>

          <c-form-block-row :label="t('scenario.workflowInCaseOfSuccess')">
            <c-workflow-field
              :model-value="parameters.stop_on_success"
              :stop-label="t('common.stop')"
              :continue-label="t('scenario.remainingStep')"
              :name="`${parametersName}.stop_on_success`"
              @update:model-value="updateParametersField('stop_on_success', $event)"
            />
          </c-form-block-row>
        </template>

        <c-form-block-row :label="t('common.comment')">
          <c-textarea
            :model-value="modelValue.comment"
            :label="t('common.comment')"
            :name="`${name}.comment`"
            @update:model-value="updateField('comment', $event)"
          />
        </c-form-block-row>
      </c-form-block>

      <c-form-general-patterns-tabs
        :form="modelValue"
        :hide-general="isPbehaviorRemoveAction"
        :with-template-testing="false"
        :patterns-label="t('common.pattern', 2)"
      >
        <template #general="{ setRef }">
          <c-form-block>
            <action-parameters-form
              v-model="parameters"
              :ref="setRef"
              :type="modelValue.type"
              :name="parametersName"
              :depth="actionParametersDepth"
              :has-previous-webhook="hasPreviousWebhook"
              :template-vars="templateVars"
            />
          </c-form-block>
        </template>

        <template #patterns="{ setRef }">
          <scenario-action-patterns-form
            :ref="setRef"
            :model-value="modelValue.patterns"
            :name="`${name}.patterns`"
            @update:model-value="updateField('patterns', $event)"
          />
        </template>
      </c-form-general-patterns-tabs>
    </c-layout>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import {
  isAssociateTicketActionType,
  isPbehaviorRemoveActionType,
  isWebhookActionType,
} from '@/helpers/entities/action/form';

import { useConfirmableForm } from '@/hooks/confirmable-form';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';
import ActionAuthorField from '@/components/other/action/form/fields/action-author-field.vue';

import ScenarioActionPatternsForm from '../scenario-action-patterns-form.vue';

export default {
  components: {
    ActionAuthorField,
    ActionParametersForm,
    ScenarioActionPatternsForm,
  },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    actionNumber: {
      type: Number,
      required: true,
    },
    hasPreviousWebhook: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    dragHandleClass: {
      type: String,
      default: 'item-drag-handler',
    },
  },

  emits: ['update:modelValue', 'remove', 'duplicate'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const parametersName = computed(() => `${props.name}.parameters.${props.modelValue.type}`);
    const parameters = computed({
      get: () => props.modelValue.parameters?.[props.modelValue.type] ?? {},
      set: value => updateField(`parameters.${props.modelValue.type}`, value),
    });
    const isWebhookAction = computed(() => isWebhookActionType(props.modelValue.type));
    const isPbehaviorRemoveAction = computed(() => isPbehaviorRemoveActionType(props.modelValue.type),
    );
    const actionParametersDepth = computed(() => (
      isAssociateTicketActionType(props.modelValue.type) ? 1 : 0
    ));

    const updateParametersField = (field, value) => {
      updateField(`parameters.${props.modelValue.type}.${field}`, value);
    };

    const { confirmAction: remove } = useConfirmableForm({
      form: toRef(props, 'modelValue'),
      action: () => emit('remove'),
    });

    const duplicate = () => emit('duplicate');

    return {
      t,
      parametersName,
      parameters,
      actionParametersDepth,
      isWebhookAction,
      isPbehaviorRemoveAction,
      updateField,
      updateParametersField,
      remove,
      duplicate,
    };
  },
};
</script>

<style lang="scss" scoped>
.scenario-action-field__actions {
  flex: 0 0 auto;
  min-height: 48px;
}
</style>
