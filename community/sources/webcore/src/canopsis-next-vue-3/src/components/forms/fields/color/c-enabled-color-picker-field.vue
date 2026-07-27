<template>
  <c-layout :column="column" :gap="column ? '2' : undefined" align-center>
    <div :class="{ 'w-100': column }" class="c-enabled-color-picker-field__switch">
      <c-enabled-field
        v-model="enabled"
        :label="label"
        color="primary"
        hide-details
        @update:model-value="updateColor"
      />
    </div>

    <v-fade-transition>
      <c-color-picker-field
        v-if="enabled"
        :model-value="modelValue"
        :disabled="!enabled"
        :required="enabled"
        :name="name"
        :justify-end="justifyEnd"
        @update:model-value="$emit('update:modelValue', $event)"
      />
    </v-fade-transition>
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
      required: false,
    },
    name: {
      type: String,
      default: 'color',
    },
    column: {
      type: Boolean,
      default: false,
    },
    justifyEnd: {
      type: Boolean,
      default: false,
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
    updateColor(value) {
      if (!value) {
        this.$emit('update:modelValue', '');
      }
    },
  },
};
</script>
