<template lang="pug">
  v-list.pa-0(dense)
    v-list-tile(
      v-for="variable in variables",
      :key="variable.value",
      :value="isActiveVariable(variable)",
      @click="selectVariable(variable)",
      @mouseenter="handleMouseEnter(variable, $event)"
    )
      v-list-tile-content
        v-list-tile-title
          v-layout(row, justify-space-between)
            | {{ variable.text }}
            span.ml-4.grey--text.lighten-1(v-if="showValue") {{ variable.value }}

      v-list-tile-action(v-if="variable.variables")
        v-icon arrow_right
    v-menu(
      v-if="subVariablesShown",
      :value="subVariablesShown",
      :position-x="subVariablesPosition.x",
      :position-y="subVariablesPosition.y",
      :z-index="zIndex",
      offset-x,
      right
    )
      variables-list(
        :variables="parentVariable.variables",
        :value="subVariableValue",
        :z-index="zIndex + 1",
        @input="selectSubVariable(parentVariable, $event)"
      )
</template>

<script>
export default {
  name: 'VariablesList',
  props: {
    value: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      default: () => [],
    },
    zIndex: {
      type: Number,
      required: false,
    },
    showValue: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      subVariablesShown: false,
      parentVariable: undefined,
      subVariablesPosition: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    subVariableValue() {
      return this.value.replace(`${this.parentVariable.value}.`, '');
    },
  },
  methods: {
    selectVariable(variable) {
      this.$emit('input', variable.value);
    },

    selectSubVariable(parentVariable, value) {
      this.$emit('input', `${parentVariable.value}.${value}`);
    },

    isActiveVariable(variable) {
      if (this.value.length > variable.value.length) {
        return this.value.startsWith(`${variable.value}.`);
      }

      return this.value.startsWith(variable.value);
    },

    handleMouseEnter(variable, event) {
      if (variable.variables) {
        const { left, top, width } = event.target.getBoundingClientRect();

        this.subVariablesPosition.x = left + width;
        this.subVariablesPosition.y = top;
        this.parentVariable = variable;
        this.subVariablesShown = true;
      } else {
        this.parentVariable = undefined;
        this.subVariablesShown = false;
      }
    },
  },
};
</script>
