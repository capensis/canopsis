<template lang="pug">
  v-layout(column)
    v-layout(row)
      div.subheading.font-weight-bold {{ title }}
      c-help-icon.ml-2(v-if="helpText", :text="helpText", max-width="300", top)
    v-layout.mt-1(row)
      span.grey--text.text--darken-2
        slot(name="subtitle")
    slot
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: false,
    },
  },
};
</script>
