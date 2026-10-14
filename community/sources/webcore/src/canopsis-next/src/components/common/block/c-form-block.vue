<template>
  <v-row class="c-form-block" no-gutters>
    <v-col v-if="title || helpText" class="c-form-block__header mb-3 d-flex align-center" cols="12">
      <div class="text-subtitle-1 font-weight-bold">
        {{ title }}
      </div>
    </v-col>

    <v-col v-if="$slots.subtitle" class="c-form-block__subtitle mb-2" cols="12">
      <span class="text-medium-emphasis">
        <slot name="subtitle" />
      </span>
    </v-col>

    <v-col :class="surfaceClasses" class="c-form-block__surface" cols="12">
      <slot />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';
import { useTheme } from 'vuetify';

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

  setup() {
    const theme = useTheme();

    const surfaceClasses = computed(() => ({
      'c-form-block__surface--dark': theme.current.value.dark,
      'c-form-block__surface--light': !theme.current.value.dark,
    }));

    return {
      surfaceClasses,
    };
  },
};
</script>

<style lang="scss">
.c-form-block__surface {
  border-radius: 8px;
  overflow: hidden;
}

.theme--light .c-form-block__surface,
.c-form-block__surface--light {
  border: 1px solid var(--v-application-background-darken2);
}

.theme--dark .c-form-block__surface,
.c-form-block__surface--dark {
  border: 1px solid var(--v-application-background-lighten3);
}
</style>
