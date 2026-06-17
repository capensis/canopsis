<template>
  <c-select-field
    v-field="value"
    :items="items"
    :label="$t('llm.model')"
    :name="name"
    :disabled="disabled"
    :loading="loading"
    :required="required"
    item-text="name"
    item-value="name"
  >
    <template #item="{ item, attrs, on }">
      <v-list-item
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-content>{{ item.name }}</v-list-item-content>
        <span
          v-if="item.recommended"
          class="ml-4 grey--text"
        >
          {{ $t('llm.recommendedBadge') }}
        </span>
      </v-list-item>
    </template>
  </c-select-field>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    loading: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'model',
    },
    required: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
