<template lang="pug">
  v-card.scenario-action-field
    v-card-text
      v-layout(row, align-center)
        v-layout.scenario-action-field__actions
          c-draggable-step-number(
            :color="hasChildrenError ? 'error' : 'primary'",
            drag-class="action-drag-handler"
          ) {{ actionNumber }}
          c-expand-btn(v-model="expanded")
        c-action-type-field.px-2(v-field="action.type", :name="`${name}.type`")
        c-action-btn(type="delete", @click="removeAction")
      v-expand-transition(mode="out-in")
        v-layout(v-show="expanded", column)
          c-enabled-field(v-field="action.emit_trigger", :label="$t('scenario.emitTrigger')")
          action-author-field(v-if="!isPbehaviorAction", v-model="parameters")
          c-workflow-field(
            v-field="action.drop_scenario_if_not_matched",
            :label="$t('scenario.workflow')",
            :continue-label="$t('scenario.remainingAction')"
          )
          v-textarea.mt-2(v-field="action.comment", :label="$tc('common.comment')")

          v-tabs(v-model="activeTab", centered, slider-color="primary", color="transparent", fixed-tabs)
            v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('common.general') }}
            v-tab(:class="{ 'error--text': hasPatternsError }") {{ $tc('common.pattern') }}
          v-divider
          v-tabs-items.pt-2(v-model="activeTab")
            v-tab-item
              action-parameters-form.mt-4(
                ref="general",
                v-model="parameters",
                :name="`${name}.parameters`",
                :type="action.type"
              )
            v-tab-item
              scenario-action-patterns-form.mt-4(
                ref="patterns",
                v-model="action.patterns",
                :name="name"
              )
</template>

<script>
import { isPbehaviorActionType } from '@/helpers/forms/action';

import { formMixin, validationChildrenMixin } from '@/mixins/form';
import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';
import ActionAuthorField from '@/components/other/action/form/partials/action-author-field.vue';

import ScenarioActionPatternsForm from './scenario-action-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ActionAuthorField,
    ActionParametersForm,
    ScenarioActionPatternsForm,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
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
    isPbehaviorAction() {
      return isPbehaviorActionType(this.action.type);
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
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
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
  },
};
</script>

<style lang="scss">
.scenario-action-field {
  &__actions {
    max-width: 100px;
  }
}
</style>
