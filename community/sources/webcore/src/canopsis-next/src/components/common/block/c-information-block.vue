<template lang="pug">
  v-layout.c-information-block(column)
    v-layout(v-if="title || helpText", row)
      div.subheading.font-weight-bold {{ title }}
      c-help-icon(
        v-if="helpText",
        :text="helpText",
        :icon="helpIcon",
        :color="helpIconColor",
        icon-class="ml-2",
        max-width="300",
        top
      )
    v-layout.mt-2(v-if="$slots.subtitle", row)
      span.text--secondary
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
    helpIcon: {
      type: String,
      required: false,
    },
    helpIconColor: {
      type: String,
      required: false,
    },
  },
};
</script>
