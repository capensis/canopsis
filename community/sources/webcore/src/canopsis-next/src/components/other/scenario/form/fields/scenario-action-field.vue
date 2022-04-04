<template lang="pug">
  v-card
    v-card-text
      v-layout(row, align-center)
        v-flex(xs2)
          v-layout
            c-draggable-step-number(
              :color="hasChildrenError ? 'error' : 'primary'",
              drag-class="action-drag-handler"
            ) {{ actionNumber }}
            c-expand-btn(v-model="expanded")
        v-flex.px-2(xs9)
          c-action-type-field(v-field="action.type", :name="`${name}.type`")
        v-flex(xs1)
          c-action-btn(type="delete", @click="removeAction")
      v-expand-transition(mode="out-in")
        v-layout(v-show="expanded", column)
          v-layout(row)
            c-enabled-field(v-field="action.emit_trigger", :label="$t('scenario.emitTrigger')")
          v-layout(row)
            c-workflow-field(
              v-field="action.drop_scenario_if_not_matched",
              :label="$t('scenario.workflow')",
              :continue-label="$t('scenario.remainingAction')"
            )
          v-layout.mt-1(row)
            v-textarea(v-field="action.comment", :label="$tc('common.comment')")
          v-tabs(v-model="activeTab", centered, slider-color="primary", color="transparent", fixed-tabs)
            v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('common.general') }}
            v-tab(:class="{ 'error--text': hasPatternsError }") {{ $tc('common.pattern') }}
          v-divider
          v-tabs-items.pt-2(v-model="activeTab")
            v-tab-item
              action-parameters-form.mt-4(
                ref="general",
                v-field="action",
                :name="`${name}.parameters`"
              )
            v-tab-item
              scenario-action-patterns-form(
                ref="patterns",
                v-model="action.patterns",
                :name="name"
              )
</template>

<script>
import { formMixin, validationChildrenMixin } from '@/mixins/form';
import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';

import ScenarioActionPatternsForm from './scenario-action-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
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
