<template lang="pug">
  v-combobox(
    ref="field",
    v-validate="rules",
    :value="value",
    :search-input="value",
    :label="label || $t('common.payload')",
    :items="availableVariables",
    :disabled="disabled",
    :return-object="false",
    :menu-props="{ value: !!variables.length && variablesShown }",
    :error-messages="errorMessages",
    :clearable="clearable",
    :name="name",
    no-filter,
    @blur="handleBlur",
    @update:searchInput="onSearchInputChange"
  )
    template(#append="")
      slot(name="append")
    template(#item="{ item, tile }")
      v-list-tile(
        v-bind="{ ...tile.props, value: item.value === variablesMenuValue }",
        @click="pasteVariable(item.value)"
      )
        v-list-tile-content {{ item.text }}
        span.ml-4.grey--text {{ item.value }}
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
        this.$validator.validate(this.name);
      }
    },
  },
};
</script>
