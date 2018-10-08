<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.criticityLevels')}}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Minor :
            v-flex(xs3)
              v-text-field(type="number", :value="value.minor", @input="e => updateValue(e, 'minor')")
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Major :
            v-flex(xs3)
              v-text-field(type="number", :value="value.major", @input="e => updateValue(e, 'major')")
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Critical :
            v-flex(xs3)
              v-text-field(type="number", :value="value.critical", @input="e => updateValue(e, 'critical')")
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
    },
  },
  methods: {
    updateValue(event, key) {
      const newValue = { ...this.value };
      newValue[key] = parseInt(event, 10);
      this.$emit('input', newValue);
    },
  },
};
</script>
