<template>
  <v-form
    class="widget-settings"
    :class="{ 'widget-settings--divider': divider }"
    @submit.prevent="$emit('submit')"
  >
    <v-list
      class="widget-settings__list py-0 mb-2"
      expand
    >
      <slot />
    </v-list>
    <v-btn
      class="mx-2 my-1"
      :loading="submitting"
      :disabled="submitting || errors.any()"
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
  --item-divider-border: 1px solid var(--item-divider-border-color);

  .theme--light & {
    --item-divider-border-color: rgba(0, 0, 0, 0.12);
  }

  .theme--dark & {
    --item-divider-border-color: rgba(255, 255, 255, 0.12);
  }

  &--divider {
    .v-list-group:not(:last-of-type) {
      border-bottom: var(--item-divider-border);
    }
  }

  &--divider &__list {
    border-bottom: var(--item-divider-border);
  }
}
</style>
