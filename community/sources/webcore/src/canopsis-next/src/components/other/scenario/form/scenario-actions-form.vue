<template lang="pug">
  v-layout(column)
    v-flex(v-show="!actions.length", xs12)
      v-alert(:value="true", type="info") {{ $t('scenario.emptyActions') }}
    c-draggable-list-field(
      v-field="actions",
      :group="draggableGroup",
      handle=".action-drag-handler",
      ghost-class="grey"
    )
      scenario-action-field.mb-3.lighten-2(
        v-for="(action, index) in actions",
        v-field="actions[index]",
        :name="`${name}.${action.key}`",
        :key="action.key",
        :action-number="index + 1",
        @remove="removeItemFromArray(index)"
      )
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasActionsErrors ? 'error' : 'primary'",
        outline,
        @click="addAction"
      ) {{ $t('scenario.addAction') }}
      span.error--text(v-show="hasActionsErrors") {{ $t('scenario.errors.actionRequired') }}
</template>

<script>
import { actionToForm } from '@/helpers/forms/action';

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
  },
};
</script>
