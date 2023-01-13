<template lang="pug">
  v-layout(align-center)
    v-layout(row)
      v-flex.mr-3(xs6)
        component(
          :is="itemTextComponent.is",
          v-field="item[itemText]",
          v-validate="'required'",
          v-bind="itemTextComponent.props",
          :error-messages="errors.collect(textFieldName)"
        )
      v-flex(xs6)
        component(
          :is="itemValueComponent.is",
          v-field="item[itemValue]",
          v-validate="",
          v-bind="itemValueComponent.props",
          :error-messages="errors.collect(valueFieldName)"
        )
    c-action-btn(v-if="!disabled", type="delete", @click="$emit('remove')")
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'item',
    event: 'input',
  },
  props: {
    item: {
      type: Object,
      default: () => ({}),
    },
    textLabel: {
      type: String,
      default: '',
    },
    valueLabel: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    name: {
      type: String,
      default: 'item',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hints: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    textFieldName() {
      return `${this.name}.${this.itemText}`;
    },

    valueFieldName() {
      return `${this.name}.${this.itemValue}`;
    },

    textHints() {
      return this.hints.map(({ text }) => text);
    },

    valueHints() {
      return this.hints.find(({ text }) => text === this.item[this.itemText])?.value ?? [];
    },

    itemTextComponent() {
      const props = {
        label: this.textLabel,
        disabled: this.disabled,
        name: this.textFieldName,
      };

      if (this.textHints.length) {
        props.items = this.textHints;
      }

      return {
        is: this.textHints.length ? 'v-combobox' : 'v-text-field',
        props,
      };
    },

    itemValueComponent() {
      const props = {
        label: this.valueLabel,
        disabled: this.disabled,
        name: this.valueFieldName,
      };

      if (this.valueHints.length) {
        props.items = this.valueHints;
      }

      return {
        is: this.valueHints.length ? 'v-combobox' : 'v-text-field',
        props,
      };
    },
  },
};
</script>
