<template>
  <v-menu
    bottom
    left
    @click.native.stop=""
  >
    <template #activator="{ on }">
      <v-btn
        icon
        v-on="on"
      >
        <v-icon>more_vert</v-icon>
      </v-btn>
    </template>
    <c-list
      :items="actions"
      item-text="title"
      item-value="title"
      return-object
      @input="selectAction"
    >
      <template #item-title="{ item }">
        <span class="mr-4">
          <v-progress-circular
            v-if="item.loading"
            :color="item.iconColor"
            :size="16"
            :width="2"
            indeterminate
          />
          <v-icon
            v-else
            :color="item.iconColor"
            :disabled="item.disabled"
            class="ma-0 pa-0"
            left
            small
          >
            {{ item.icon }}
          </v-icon>
        </span>
        <span
          :class="item.cssClass"
          class="text-body-1"
        >
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
  },
  setup() {
    const selectAction = item => item?.method?.();

    return {
      selectAction,
    };
  },
};
</script>
