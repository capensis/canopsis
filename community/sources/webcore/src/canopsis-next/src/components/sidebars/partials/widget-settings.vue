<template>
  <v-form
    :class="{ 'widget-settings--divider': divider }"
    class="widget-settings"
    @submit.prevent="$emit('submit')"
  >
    <v-list
      class="widget-settings__list py-0"
      expand
    >
      <slot />
    </v-list>
    <v-layout class="widget-settings__submit-btn-wrapper pa-4">
      <v-btn
        :loading="submitting"
        :disabled="submitting || !validatorDirty"
        type="submit"
        color="primary"
      >
        {{ $t('common.save') }}
      </v-btn>
    </v-layout>
  </v-form>
</template>

<script>
export default {
  inject: ['$validator', '$sidebar'],
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
  computed: {
    validatorDirty() {
      return !this.$sidebar?.config?.widget?._id || Object.values(this.fields).some(field => field.dirty);
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

  &__submit-btn-wrapper {
    background-color: var(--v-application-background-base);
    position: sticky;
    bottom: 0;

    .theme--light & {
      background-color: #fff;
    }
  }
}
</style>
