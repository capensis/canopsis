import { debounce } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  findSelectedVariable,
  matchPayloadOperators,
  matchPayloadVariables,
} from '@/helpers/payload-json';

export const payloadFieldMixin = {
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

  emits: ['update:value', 'blur'],

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
    variablesMenuValue() {
      return this.variableGroup?.[0] ?? '';
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

  beforeUnmount() {
    document.removeEventListener('selectionchange', this.debouncedOnSelectionChange);
  },

  methods: {
    getInputElement() {
      return this.$refs.field?.$el?.querySelector('input, textarea');
    },

    getVariableElement() {
      const variable = this.$refs.variable;

      return Array.isArray(variable) ? variable[0] : variable;
    },

    updateModel(value) {
      this.$emit('update:value', value ?? '');
      this.$nextTick(this.onSelectionChange);
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
      const cursorPosition = this.selectionVariableStart + variable.length;
      const nextValue = `${this.value.slice(0, this.selectionVariableStart)}${variable}${this.value.slice(this.selectionVariableEnd)}`;

      this.updateModel(nextValue);
      this.resetVariableSelection();
      this.hideVariablesMenu();

      this.$nextTick(() => {
        const input = this.getInputElement();

        if (!input) {
          return;
        }

        input.focus();
        input.setSelectionRange(cursorPosition, cursorPosition);
      });
    },

    showVariablesMenu() {
      this.variablesShown = true;

      this.$nextTick(() => {
        const variableElement = this.getVariableElement();

        if (!variableElement) {
          return;
        }

        const { top, left, height } = variableElement.getBoundingClientRect();

        this.variablesMenuPosition.x = left;
        this.variablesMenuPosition.y = top + height;
      });
    },

    hideVariablesMenu() {
      this.variablesShown = false;
    },

    isVariableCreatingInsideOperatorContent() {
      const { index, groups } = this.operatorGroup;

      return (
        this.newVariableGroup.index > index + groups.open.length &&
        this.newVariableGroup.index < index + groups.variable.length - groups.close.length
      );
    },

    onSelectionChange() {
      const input = this.getInputElement();

      if (!input || !this.$el.contains(document.activeElement)) {
        return;
      }

      const { selectionStart, selectionEnd } = input;

      this.variableGroup = findSelectedVariable(this.valueVariables, selectionStart, selectionEnd);
      this.operatorGroup = findSelectedVariable(this.valueOperators, selectionStart, selectionEnd);
      this.newVariableGroup = this.value.slice(0, selectionStart).match(/({{){1,2}$/);

      if (this.newVariableGroup && !this.variableGroup) {
        if (!this.operatorGroup || this.isVariableCreatingInsideOperatorContent()) {
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

    handleBlur(event) {
      setTimeout(() => {
        this.hideVariablesMenu();
        this.resetVariableSelection();
      }, VUETIFY_ANIMATION_DELAY);

      this.$emit('blur', event);
    },
  },
};

export default payloadFieldMixin;
