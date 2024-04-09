import { debounce } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { findSelectedVariable, matchPayloadOperators, matchPayloadVariables } from '@/helpers/payload-json';

import { formBaseMixin } from '@/mixins/form';

export const payloadFieldMixin = {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectionVariableStart: 0,
      selectionVariableEnd: 0,
      variableGroup: undefined,
      operatorGroup: undefined,
      newVariableGroup: undefined,
      variablesShown: false,
      variablesMenuPosition: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    availableVariables() {
      return this.prepareVariables(this.variables);
    },

    variablesMenuValue() {
      return this.variableGroup?.[0] && this.operatorGroup?.[0];
    },

    valueVariables() {
      return matchPayloadVariables(this.value);
    },

    valueOperators() {
      return matchPayloadOperators(this.value);
    },
  },
  created() {
    this.debouncedOnSelectionChange = debounce(this.onSelectionChange, 300);
  },
  mounted() {
    document.addEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  beforeDestroy() {
    document.removeEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  methods: {
    prepareSubVariables(variables, parentVariable) {
      return variables.map(variable => this.prepareVariable(variable, parentVariable));
    },

    prepareVariable(variable, parentVariable) {
      return {
        ...variable,
        variables: variable.variables ? this.prepareSubVariables(variable.variables, variable) : variable.variables,
        value: `{{ ${parentVariable ? parentVariable.value : ''}${variable.value} }}`,
      };
    },

    prepareVariableWithEnumerable(variable) {
      return variable.variables.map(subVariable => ({
        ...this.prepareVariable(subVariable, variable),
        value: (this.variableGroup || this.newVariableGroup) && this.operatorGroup
          ? `{{ ${subVariable.value} }}`
          : `{{ range ${variable.value} }}{{ ${subVariable.value} }}{{ end }}`,
      }));
    },

    prepareVariables(variables) {
      return variables.reduce((acc, variable) => {
        if (variable.enumerable) {
          acc.push(...this.prepareVariableWithEnumerable(variable));
        } else {
          acc.push(this.prepareVariable(variable));
        }

        return acc;
      }, []);
    },

    setVariableSelection(start, end) {
      this.selectionVariableStart = start;
      this.selectionVariableEnd = end;
    },

    setVariableSelectionByGroup(group) {
      const [value] = group;

      this.setVariableSelection(group.index, group.index + value.length);
    },

    resetVariableSelection() {
      this.selectionVariableStart = undefined;
      this.selectionVariableEnd = undefined;
    },

    pasteVariable(variable) {
      const newValue = `${this.value.slice(0, this.selectionVariableStart)}${variable}${this.value.slice(this.selectionVariableEnd)}`;

      this.updateModel(newValue);
      this.resetVariableSelection();
      this.hideVariablesMenu();

      if (this.errors && this.name) {
        this.errors.remove(this.name);
      }
    },

    showVariablesMenu() {
      this.variablesShown = true;

      this.$nextTick(() => {
        if (this.$refs.variable) {
          const [variableElement] = this.$refs.variable;
          const { top, left, height } = variableElement.getBoundingClientRect();

          this.variablesMenuPosition.x = left;
          this.variablesMenuPosition.y = top + height;
        }
      });
    },

    hideVariablesMenu() {
      this.variablesShown = false;
    },

    isVariableCreatingInsideOperatorContent() {
      const { index, groups } = this.operatorGroup;

      return (
        this.newVariableGroup.index > index + groups.open.length
        && this.newVariableGroup.index < index + groups.variable.length - groups.close.length
      );
    },

    onSelectionChange() {
      if (!this.$el.contains(document.activeElement)) {
        return;
      }

      const { selectionStart, selectionEnd } = this.$refs.field.$refs.input;

      this.variableGroup = findSelectedVariable(
        this.valueVariables,
        selectionStart,
        selectionEnd,
      );
      this.operatorGroup = findSelectedVariable(
        this.valueOperators,
        selectionStart,
        selectionEnd,
      );
      this.newVariableGroup = this.value.slice(0, selectionStart).match(/({{){1,2}$/);

      if (this.newVariableGroup && !this.variableGroup) {
        if (
          !this.operatorGroup
          || this.isVariableCreatingInsideOperatorContent()
        ) {
          this.setVariableSelection(this.newVariableGroup.index, selectionEnd);
          this.showVariablesMenu();
          return;
        }
      }

      if (this.variableGroup || this.operatorGroup) {
        this.setVariableSelectionByGroup(this.variableGroup || this.operatorGroup);
        this.showVariablesMenu();
        return;
      }

      this.hideVariablesMenu();
      this.resetVariableSelection();
    },

    handleBlur() {
      setTimeout(() => {
        this.hideVariablesMenu();
        this.resetVariableSelection();
      }, VUETIFY_ANIMATION_DELAY);
    },
  },
};
