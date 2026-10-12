<template>
  <c-layout column>
    <c-enabled-field
      v-model="enabled"
      :label="label"
      color="primary"
      @update:model-value="updateIcon"
    />
    <c-icon-field
      v-if="enabled"
      :model-value="modelValue"
      :label="$t('common.icon', 1)"
      :name="name"
      required
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon',
    },
  },

  emits: ['update:modelValue'],

  data() {
    return {
      enabled: !!this.modelValue,
    };
  },

  watch: {
    modelValue(value) {
      this.enabled = !!value;
    },
  },

  methods: {
    updateIcon(value) {
      if (!value) {
        this.$emit('update:modelValue', '');
      }
    },
  },
};
</script>
