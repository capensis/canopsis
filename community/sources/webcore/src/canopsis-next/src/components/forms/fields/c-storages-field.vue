<template>
  <v-row>
    <v-col v-show="label" cols="12">
      <c-layout>
        <h4 class="text-subtitle-1 text-grey-darken-2">
          {{ label }}
        </h4>
        <c-help-icon
          v-if="helpText"
          :text="helpText"
          icon-class="ml-2 storage-help-tooltip"
          right
        />
      </c-layout>
    </v-col>
    <v-col cols="12">
      <c-layout column>
        <c-storage-field
          v-for="(storage, index) in modelValue"
          :key="storage.key"
          :model-value="storage.directory"
          :disabled="disabled"
          @edit="$emit('edit', storage, index)"
          @remove="removeItemFromArray(index)"
        />
      </c-layout>
    </v-col>
    <v-col cols="12">
      <c-layout>
        <v-btn
          :disabled="disabled"
          class="ml-0"
          color="primary"
          @click="$emit('add')"
        >
          {{ $t('common.add') }}
        </v-btn>
      </c-layout>
    </v-col>
  </v-row>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue', 'edit', 'add'],

  methods: {
    removeItemFromArray(index) {
      this.$emit('update:modelValue', this.modelValue.filter((v, i) => i !== index));
    },
  },
};
</script>

<style lang="scss">
.storage-help-tooltip {
  pointer-events: auto !important;
}
</style>
