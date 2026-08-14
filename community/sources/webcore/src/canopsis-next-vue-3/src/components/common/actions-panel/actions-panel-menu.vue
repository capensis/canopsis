<template>
  <v-menu v-if="actions.length" location="bottom end" @click.stop="">
    <template #activator="{ props }">
      <c-action-btn
        :tooltip="$t('common.moreButton')"
        :small="small"
        icon="more_vert"
        v-bind="props"
      />
    </template>

    <c-list
      :items="actions"
      item-text="title"
      item-value="title"
      return-object
      @update:value="selectAction"
    >
      <template #item-title="{ item }">
        <span class="actions-panel__menu-item-icon mr-4">
          <v-progress-circular
            v-if="item.loading"
            :color="item.iconColor"
            :size="16"
            :width="2"
            indeterminate
          />
          <c-icon
            v-else
            :color="item.iconColor"
            :disabled="item.disabled"
            :size="item.iconSize"
            :small="!item.iconSize"
            class="ma-0 pa-0"
          >
            {{ item.icon }}
          </c-icon>
        </span>
        <span :class="item.cssClass" class="text-body-1">
          {{ item.title }}
        </span>
      </template>
    </c-list>
  </v-menu>
</template>

<script>
export default {
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    small: {
      type: Boolean,
      default: false,
    },
  },

  methods: {
    selectAction(item) {
      item?.method?.();
    },
  },
};
</script>
