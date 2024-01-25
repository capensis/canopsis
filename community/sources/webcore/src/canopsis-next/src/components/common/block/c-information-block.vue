<template>
  <v-layout
    class="c-information-block"
    column
  >
    <v-layout
      v-if="title || helpText"
    >
      <div class="text-subtitle-1 font-weight-bold">
        {{ title }}
      </div>
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        :icon="helpIcon"
        :color="helpIconColor"
        icon-class="ml-2"
        max-width="300"
        top
      />
    </v-layout>
    <v-layout
      class="mt-2"
      v-if="$slots.subtitle"
    >
      <span class="text--secondary">
        <slot name="subtitle" />
      </span>
    </v-layout>
    <slot />
  </v-layout>
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
