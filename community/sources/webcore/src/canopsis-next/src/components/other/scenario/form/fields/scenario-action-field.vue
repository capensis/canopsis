<template>
  <c-card-iterator-item
    ref="cardIteratorItemElement"
    :item-number="actionNumber"
    @remove="removeAction"
  >
    <template #header="">
      <c-action-type-field
        v-field="action.type"
        :name="`${name}.type`"
      />

      <c-action-btn type="duplicate" @click="duplicateAction" />
    </template>

    <v-layout class="gap-4" column>
      <c-form-block>
        <c-form-block-row :label="$t('common.emitTrigger')">
          <c-enabled-field
            v-field="action.emit_trigger"
            :label="$t('common.emitTrigger')"
          />
        </c-form-block-row>

        <c-form-block-row :label="$t('scenario.forwardAuthor')">
          <action-author-field v-model="parameters" :variables="templateVars.author" />
        </c-form-block-row>

        <c-form-block-row v-if="isWebhookAction" :label="$t('scenario.skip')">
          <c-enabled-field
            v-model="parameters.skip_for_child"
            :label="$t('scenario.skipForChild')"
            hide-details
          />
          <c-enabled-field
            v-model="parameters.skip_for_instruction"
            :label="$t('scenario.skipForInstruction')"
          />
        </c-form-block-row>

        <c-form-block-row :label="$t('scenario.workflow')" indented>
          <v-layout justify-space-between>
            <c-workflow-field
              v-field="action.drop_scenario_if_not_matched"
              :label="$t('scenario.workflow')"
              :continue-label="$t('scenario.remainingAction')"
            />

            <template v-if="isWebhookAction">
              <c-workflow-field
                v-model="parameters.stop_on_fail"
                :label="$t('scenario.workflowInCaseOfFailure')"
                :continue-label="$t('scenario.remainingStep')"
              />
              <c-workflow-field
                v-model="parameters.stop_on_success"
                :label="$t('scenario.workflowInCaseOfSuccess')"
                :continue-label="$t('scenario.remainingStep')"
              />
            </template>
          </v-layout>
        </c-form-block-row>

        <c-form-block-row :label="$tc('common.comment')">
          <v-textarea
            v-field="action.comment"
            :label="$tc('common.comment')"
          />
        </c-form-block-row>
      </c-form-block>

      <c-form-general-patterns-tabs :hide-general="isPbehaviorRemoveAction">
        <template #general="{ setRef }">
          <c-form-block>
            <action-parameters-form
              v-model="parameters"
              :ref="setRef"
              :name="`${name}.parameters`"
              :type="action.type"
              :has-previous-webhook="hasPreviousWebhook"
              :template-vars="templateVars"
            />
          </c-form-block>
        </template>
        <template #patterns="{ setRef }">
          <scenario-action-patterns-form
            v-field="action.patterns"
            :ref="setRef"
            :name="name"
          />
        </template>
      </c-form-general-patterns-tabs>
    </v-layout>
  </c-card-iterator-item>
</template>

<script>
import {
  computed,
  ref,
  toRef,
  onMounted,
  inject,
} from 'vue';
import { Validator } from 'vee-validate';

import { isWebhookActionType, isPbehaviorRemoveActionType } from '@/helpers/entities/action';

import { useConfirmableForm } from '@/hooks/confirmable-form';
import { useModelField } from '@/hooks/form/model-field';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';
import ActionAuthorField from '@/components/other/action/form/fields/action-author-field.vue';

import ScenarioActionPatternsForm from '../scenario-action-patterns-form.vue';

export default {
  inject: {
    $validator: {
      default: () => new Validator(),
    },
    $aiChat: {
      default: () => ({}),
    },
  },
  components: {
    ActionAuthorField,
    ActionParametersForm,
    ScenarioActionPatternsForm,
  },
  model: {
    prop: 'action',
    event: 'input',
  },
  props: {
    action: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'action',
    },
    actionNumber: {
      type: [Number, String],
      default: 0,
    },
    hasPreviousWebhook: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const $aiChat = inject('$aiChat', {});

    const { updateField } = useModelField(props, emit);

    const cardIteratorItemElement = ref(null);

    const { confirmAction: removeAction } = useConfirmableForm({
      form: toRef(props, 'action'),
      action: () => emit('remove'),
      cloning: true,
    });

    const isWebhookAction = computed(() => isWebhookActionType(props.action.type));
    const isPbehaviorRemoveAction = computed(() => isPbehaviorRemoveActionType(props.action.type));

    const parameters = computed({
      get() {
        const { type, parameters: actionParameters } = props.action;

        return actionParameters[type];
      },
      set(value) {
        updateField(`parameters.${props.action.type}`, value);
      },
    });

    const duplicateAction = () => emit('duplicate');

    const toggleOnExpanded = () => cardIteratorItemElement.value.toggleOnExpanded();

    onMounted(() => {
      $aiChat?.registerExpandFunction?.(({ key }) => {
        if (key === props.action.key) {
          toggleOnExpanded();
        }
      });
    });

    return {
      cardIteratorItemElement,

      isWebhookAction,
      isPbehaviorRemoveAction,
      parameters,
      removeAction,
      duplicateAction,
    };
  },
};
</script>
