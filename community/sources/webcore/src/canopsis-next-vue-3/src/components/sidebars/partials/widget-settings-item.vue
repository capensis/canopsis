<template>
  <v-list-group
    v-bind="$attrs"
    :eager="eager"
    class="widget-settings-item"
    color=""
  >
    <template #activator="{ props }">
      <v-list-item v-bind="props">
        <template #title>
          <c-layout align="center" gap="1">
            <span>
              <slot name="title">
                {{ title }}
              </slot>
            </span>

            <span
              v-if="optional || subTitle"
              class="font-italic text-caption text-medium-emphasis"
            >
              <span v-if="optional" class="text-lowercase">
                ({{ $t('common.optional') }})
              </span>

              <span v-if="subTitle">
                {{ subTitle }}
              </span>
            </span>
          </c-layout>
        </template>
      </v-list-item>
    </template>

    <v-container class="widget-settings-item__content">
      <slot />
    </v-container>
  </v-list-group>
</template>

<script>
export default {
  inheritAttrs: false,

  props: {
    title: {
      type: String,
      default: '',
    },
    subTitle: {
      type: String,
      default: '',
    },
    optional: {
      type: Boolean,
      default: false,
    },
    eager: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-settings-item {
  &__content {
    padding: 24px;
  }

  :deep(.v-list-item) {
    padding-top: 0;
    padding-bottom: 0;
  }

  :deep(.v-list-group__items .v-list-item) {
    padding-inline-start: 16px !important;
  }
}
</style>
