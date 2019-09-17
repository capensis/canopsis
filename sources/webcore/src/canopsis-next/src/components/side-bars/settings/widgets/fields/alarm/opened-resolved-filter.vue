<template lang="pug">
  v-list-group(data-test="filterOnOpenResolved")
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-checkbox(
          v-model="value.opened",
          :label="$t('settings.open')",
          data-test="openFilter",
          color="primary",
          hide-details,
          @change="updateField('opened', $event)"
        )
        v-checkbox(
          v-model="value.resolved",
          :label="$t('settings.resolved')",
          data-test="resolvedFilter",
          color="primary",
          hide-details,
          @change="updateField('resolved', $event)"
        )
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    updateField(key, value) {
      this.$emit('input', { ...this.value, [key]: value });
    },
  },
};
</script>
