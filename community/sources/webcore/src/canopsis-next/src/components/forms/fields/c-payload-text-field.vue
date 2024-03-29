<template>
  <v-combobox
    v-validate="rules"
    ref="field"
    :value="value"
    :search-input="value"
    :label="label || $t('common.payload')"
    :items="availableVariables"
    :disabled="disabled"
    :return-object="false"
    :menu-props="{ value: !!variables.length && variablesShown, offsetY: true }"
    :error-messages="errorMessages"
    :clearable="clearable"
    :name="name"
    no-filter
    @blur="handleBlur"
    @update:search-input="onSearchInputChange"
  >
    <template #append="">
      <slot name="append" />
    </template>
    <template #item="{ item, attrs }">
      <v-list-item
        v-bind="{ ...attrs, value: item.value === variablesMenuValue }"
        @click="pasteVariable(item.value)"
      >
        <v-list-item-content>{{ item.text }}</v-list-item-content>
        <span class="ml-4 grey--text">{{ item.value }}</span>
      </v-list-item>
    </template>
  </v-combobox>
</template>

<script>
import { payloadFieldMixin } from '@/mixins/payload/payload-field';

export default {
  inject: ['$validator'],
  mixins: [payloadFieldMixin],
  props: {
    name: {
      type: String,
      default: 'payload',
    },
    label: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
    errorMessages() {
      return this.errors.collect(this.name);
    },
  },
  methods: {
    onSearchInputChange(value) {
      this.debouncedOnSelectionChange();
      this.updateModel(value ?? '');

      if (this.errorMessages?.length) {
        this.$nextTick(() => {
          this.$validator.validate(this.name);
        });
      }
    },
  },
};
</script>
