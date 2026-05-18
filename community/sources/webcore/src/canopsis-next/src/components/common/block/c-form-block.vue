<template>
  <v-layout
    class="c-form-block"
    column
  >
    <v-layout
      v-if="title || helpText"
      class="c-form-block__header mb-3"
      align-center
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
      v-if="$slots.subtitle"
      class="c-form-block__subtitle mb-2"
    >
      <span class="text--secondary">
        <slot name="subtitle" />
      </span>
    </v-layout>
    <div class="c-form-block__surface">
      <slot />
    </div>
  </v-layout>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
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

<style lang="scss">
.c-form-block__surface {
  border-radius: 8px;
  overflow: hidden;
}

.theme--light .c-form-block__surface {
  border: 1px solid var(--v-application-background-darken2);
}

.theme--dark .c-form-block__surface {
  border: 1px solid var(--v-application-background-lighten3);
}
</style>
