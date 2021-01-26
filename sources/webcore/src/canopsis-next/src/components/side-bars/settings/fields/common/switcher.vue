<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      v-flex(xs11)
        div.subheading {{ title }}
      v-layout(data-test="switcherLayout")
        v-switch.ma-0.pa-0(
          data-test="switcherField",
          color="primary",
          :input-value="value",
          hide-details,
          @change="$emit('input', $event)"
        )
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
