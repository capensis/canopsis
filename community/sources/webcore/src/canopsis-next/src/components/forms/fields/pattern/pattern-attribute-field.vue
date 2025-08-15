<template>
  <c-select-field
    v-field="value"
    v-validate="rules"
    :items="items"
    :disabled="disabled"
    :error-messages="errors.collect(name)"
    :label="label || $tc('common.searchBy')"
    :return-object="returnObject"
    :name="name"
    item-disabled="options.disabled"
    ellipsis
    autocomplete
  >
    <template #item="{ item }">
      <span>{{ item.text }}</span>
      <v-tooltip v-if="item.options?.alias" offset-y top>
        <template #activator="{ on }">
          <v-icon
            class="ml-1"
            color="primary"
            small
            v-on="on"
          >
            alternate_email
          </v-icon>
        </template>
        <span>infos.{{ item.originalValue }}.value</span>
      </v-tooltip>
    </template>
    <template #selection="{ item }">
      <v-icon
        v-if="item.options?.alias"
        class="mr-1"
        small
      >
        alternate_email
      </v-icon>
      <span>{{ item.text }}</span>
    </template>
  </c-select-field>
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [String, Object],
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'attribute',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
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
  },
};
</script>
