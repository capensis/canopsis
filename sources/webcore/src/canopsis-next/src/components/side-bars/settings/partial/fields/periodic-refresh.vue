<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
    v-container
      v-layout
        v-flex
          v-switch(
          v-model="value.enabled",
          @change="updateField('enabled', $event)",
          color="green darken-3",
          hide-details
          )
        v-flex
          v-text-field.pt-0(
          type="number",
          :value="value.interval",
          :disabled="!value.enabled",
          @input="updateField('interval', $event)",
          hide-details
          )
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      default: () => ({ enabled: false, interval: '' }),
    },
  },
  data() {
    return {
      isEnabled: false,
    };
  },
  methods: {
    updateField(fieldKey, value) {
      this.$emit('input', { ...this.value, [fieldKey]: value });
    },
  },
};
</script>
