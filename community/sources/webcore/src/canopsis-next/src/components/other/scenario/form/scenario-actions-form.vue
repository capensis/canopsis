<template lang="pug">
  v-layout.mt-2(column)
    v-layout(v-show="!actions.length", row)
      v-flex
        v-alert(:value="true", type="info") {{ $t('scenario.emptyActions') }}
    draggable(v-field="actions", :options="draggableOptions")
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
        outline,
        :color="hasActionsErrors ? 'error' : 'primary'",
        @click="addAction"
      ) {{ $t('scenario.addAction') }}
      span.error--text(v-show="hasActionsErrors") {{ $t('scenario.errors.actionRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { scenarioActionToForm } from '@/helpers/forms/scenario';

import { formArrayMixin, validationChildrenMixin } from '@/mixins/form';

import ScenarioActionField from './fields/scenario-action-field.vue';

export default {
  inject: ['$validator'],
  components: { ScenarioActionField, Draggable },
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

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.action-drag-handler',
        ghostClass: 'grey',
        group: {
          name: 'scenarios-actions',
        },
      };
    },
  },
  watch: {
    actions() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.actions.length,
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    addAction() {
      this.addItemIntoArray(scenarioActionToForm());
    },
  },
};
</script>
