<template>
  <c-card-iterator-item
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
    <v-layout>
      <v-flex xs6>
        <c-enabled-field
          v-field="action.emit_trigger"
          :label="$t('common.emitTrigger')"
        />
        <action-author-field v-model="parameters" />
      </v-flex>
      <v-flex
        v-if="isWebhookActionType"
        xs6
      >
        <c-enabled-field
          v-model="parameters.skip_for_child"
          :label="$t('scenario.skipForChild')"
        />
        <c-enabled-field
          v-model="parameters.skip_for_instruction"
          :label="$t('scenario.skipForInstruction')"
          class="mt-0"
        />
      </v-flex>
    </v-layout>
    <c-workflow-field
      v-field="action.drop_scenario_if_not_matched"
      :label="$t('scenario.workflow')"
      :continue-label="$t('scenario.remainingAction')"
    />
    <v-textarea
      v-field="action.comment"
      :label="$tc('common.comment')"
      class="mt-2"
    />
    <v-tabs
      v-model="activeTab"
      slider-color="primary"
      background-color="transparent"
      centered
    >
      <v-tab v-if="!isPbehaviorRemoveAction" :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasPatternsError }">
        {{ $tc('common.pattern') }}
      </v-tab>
    </v-tabs>
    <v-divider />
    <v-tabs-items
      v-model="activeTab"
      class="pt-2"
    >
      <v-tab-item v-if="!isPbehaviorRemoveAction" eager>
        <action-parameters-form
          v-model="parameters"
          ref="general"
          :name="`${name}.parameters`"
          :type="action.type"
          :has-previous-webhook="hasPreviousWebhook"
          class="mt-4"
        />
      </v-tab-item>
      <v-tab-item eager>
        <scenario-action-patterns-form
          v-field="action.patterns"
          ref="patterns"
          :name="name"
          class="mt-4"
        />
      </v-tab-item>
    </v-tabs-items>
  </c-card-iterator-item>
</template>

<script>
import { isPbehaviorRemoveActionType, isWebhookActionType } from '@/helpers/entities/action';

import { formMixin } from '@/mixins/form';
import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';
import ActionAuthorField from '@/components/other/action/form/fields/action-author-field.vue';

import ScenarioActionPatternsForm from '../scenario-action-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ActionAuthorField,
    ActionParametersForm,
    ScenarioActionPatternsForm,
  },
  mixins: [
    formMixin,
    confirmableFormMixinCreator({
      field: 'action',
      method: 'removeAction',
      cloning: true,
    }),
  ],
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
  },
  data() {
    return {
      activeTab: 0,
      expanded: true,
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  computed: {
    isPbehaviorRemoveAction() {
      return isPbehaviorRemoveActionType(this.action.type);
    },

    isWebhookActionType() {
      return isWebhookActionType(this.action.type);
    },

    parameters: {
      get() {
        const { type, parameters } = this.action;

        return parameters[type];
      },
      set(value) {
        this.updateField(`parameters.${this.action.type}`, value);
      },
    },
  },
  mounted() {
    this.$watch(() => this.$refs.general?.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
      this.hasPatternsError = value;
    });
  },
  methods: {
    removeAction() {
      this.$emit('remove');
    },

    duplicateAction() {
      this.$emit('duplicate');
    },
  },
};
</script>
