<template>
  <v-form
    :class="{ 'widget-settings--divider': divider }"
    class="widget-settings"
    @submit.prevent="emit('submit')"
  >
    <v-list class="widget-settings__list py-0">
      <slot />
    </v-list>

    <c-layout class="widget-settings__submit-btn-wrapper pa-4">
      <v-btn
        :loading="submitting"
        :disabled="submitting"
        type="submit"
        color="primary"
      >
        {{ t('common.save') }}
      </v-btn>
    </c-layout>
  </v-form>
</template>

<script setup>
import { useI18n } from '@/hooks/i18n';

defineProps({
  submitting: {
    type: Boolean,
    default: false,
  },
  dirty: {
    type: Boolean,
    default: false,
  },
  divider: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['submit']);

const { t } = useI18n();
</script>

<style lang="scss">
.widget-settings {
  --item-divider-border: 1px solid var(--v-divider-border-color);

  &--divider {
    .v-list-group,
    .widget-settings-flat-item {
      &:not(:last-of-type) {
        border-bottom: var(--item-divider-border);
      }
    }
  }

  &--divider &__list {
    border-bottom: var(--item-divider-border);
  }

  &__submit-btn-wrapper {
    position: sticky;
    bottom: 0;
    z-index: 2;
    background-color: rgb(var(--v-theme-surface));
  }
}
</style>
