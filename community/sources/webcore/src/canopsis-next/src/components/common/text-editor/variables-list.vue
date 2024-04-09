<template>
  <v-list
    :dense="dense"
    class="pa-0"
  >
    <v-list-item
      v-for="variable in variables"
      :key="variable.value"
      :value="isActiveVariable(variable)"
      @click="selectVariable(variable)"
      @mouseenter="handleMouseEnter(variable, $event)"
    >
      <v-list-item-content>
        <v-list-item-title>
          <v-layout justify-space-between>
            {{ variable.text }}<span
              v-if="showValue"
              class="ml-4 grey--text lighten-1"
            >{{ variable.value }}</span>
          </v-layout>
        </v-list-item-title>
      </v-list-item-content>
      <v-list-item-action v-if="variable.variables">
        <v-icon>arrow_right</v-icon>
      </v-list-item-action>
    </v-list-item>
    <v-menu
      v-if="subVariablesShown"
      :value="subVariablesShown"
      :position-x="subVariablesPosition.x"
      :position-y="subVariablesPosition.y"
      :z-index="zIndex"
      offset-x
      right
    >
      <variables-list
        :variables="parentVariable.variables"
        :value="subVariableValue"
        :z-index="zIndex + 1"
        :show-value="showValue"
        @input="selectSubVariable"
      />
    </v-menu>
  </v-list>
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
    dense: {
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

    selectSubVariable(value) {
      this.$emit('input', value);
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
