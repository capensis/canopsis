<template lang="pug">
  v-list-group(data-test="filterOnOpenResolved")
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-checkbox(
        data-test="openFilter",
        color="primary",
        :label="$t('settings.open')",
        v-model="value.opened",
        @change="updateField('opened', $event)"
        hide-details
        )
        v-checkbox(
        data-test="resolvedFilter",
        color="primary",
        :label="$t('settings.resolved')",
        v-model="value.resolved",
        @change="updateField('resolved', $event)"
        hide-details
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
