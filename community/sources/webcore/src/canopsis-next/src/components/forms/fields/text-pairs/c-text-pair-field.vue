<template>
  <v-layout class="gap-3" align-center>
    <v-text-field
      v-field="item[itemText]"
      v-validate="textValidationRules"
      :label="textLabel"
      :disabled="disabled"
      :name="textFieldName"
      :error-messages="errors.collect(textFieldName)"
    />

    <c-select-field
      v-if="items.length"
      v-field="item[itemValue]"
      :label="valueLabel"
      :disabled="disabled"
      :items="items"
      :name="valueFieldName"
      :required="valueRequired"
    />
    <c-payload-text-field
      v-else
      v-field="item[itemValue]"
      :label="valueLabel"
      :disabled="disabled"
      :name="valueFieldName"
      :variables="variables"
      :required="valueRequired"
    >
      <template #append="">
        <slot name="append-value" />
      </template>
    </c-payload-text-field>

    <c-action-btn
      v-if="!disabled"
      type="delete"
      @click="$emit('remove')"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

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
    textRequired: {
      type: Boolean,
      default: false,
    },
    valueRequired: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const textValidationRules = computed(() => ({
      required: props.textRequired,
    }));

    const textFieldName = computed(() => `${props.name}.${props.itemText}`);
    const valueFieldName = computed(() => `${props.name}.${props.itemValue}`);

    return {
      textValidationRules,
      textFieldName,
      valueFieldName,
    };
  },
};
</script>
