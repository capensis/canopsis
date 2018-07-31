<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-checkbox(
        :label="$t('settings.open')",
        v-model="value.opened",
        @change="change('opened')($event)"
        hide-details
        )
        v-checkbox(
        :label="$t('settings.resolved')",
        v-model="value.resolved",
        @change="change('resolved')($event)"
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
    change(key) {
      return value => this.$emit('input', { ...this.value, [key]: value });
    },
  },
};
</script>
