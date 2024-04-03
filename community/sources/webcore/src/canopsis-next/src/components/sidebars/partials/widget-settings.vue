<template>
  <v-form
    :class="{ 'widget-settings--divider': divider }"
    class="widget-settings"
    @submit.prevent="$emit('submit')"
  >
    <v-list
      class="widget-settings__list py-0 mb-2"
      expand
    >
      <slot />
    </v-list>
    <v-btn
      :loading="submitting"
      :disabled="submitting || errors.any()"
      class="mx-2 my-1"
      type="submit"
      color="primary"
    >
      {{ $t('common.save') }}
    </v-btn>
  </v-form>
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    submitting: {
      type: Boolean,
      default: false,
    },
    divider: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.widget-settings {
  --item-divider-border: 1px solid var(--v-divider-border-color);

  &--divider {
    .v-list-group, .widget-settings-flat-item {
      &:not(:last-of-type) {
        border-bottom: var(--item-divider-border);
      }
    }
  }

  &--divider &__list {
    border-bottom: var(--item-divider-border);
  }
}
</style>
