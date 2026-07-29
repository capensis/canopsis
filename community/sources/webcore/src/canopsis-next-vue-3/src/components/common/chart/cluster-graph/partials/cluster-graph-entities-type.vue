<template>
  <div>
    <h6 class="my-2 text-h6 text-center">
      {{ $t('common.type') }}
    </h6>
    <c-layout justify-center>
      <v-radio-group
        :model-value="impact"
        class="my-0"
        inline
        hide-details
        @update:model-value="updateValue"
      >
        <v-radio :value="false" :label="$t('context.treeOfDependencies')" color="primary" />
        <v-radio :value="true" :label="$t('context.impactChain')" color="primary" />
      </v-radio-group>
    </c-layout>
  </div>
</template>

<script>
export default {
  props: {
    impact: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:impact', 'change'],
  methods: {
    updateValue(value) {
      this.$emit('update:impact', value);
      this.$emit('change', value);
    },
  },
};
</script>
