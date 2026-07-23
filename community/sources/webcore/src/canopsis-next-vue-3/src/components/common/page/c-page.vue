<template>
  <div>
    <c-page-header :name="name">
      <slot v-if="$slots.header" name="header" />
    </c-page-header>

    <v-card class="mx-4 mb-4 mt-0">
      <slot />
    </v-card>

    <c-fab-btn
      :has-access="creatable"
      :has-create-listener="hasCreateListener"
      :dark="dark"
      @refresh="$emit('refresh')"
      @create="$emit('create')"
    >
      <span>{{ createTooltip }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      default: undefined,
    },
    creatable: {
      type: Boolean,
      default: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    createTooltip: {
      type: String,
      required: false,
    },
    hasCreateListener: {
      type: Boolean,
      default: true,
    },
  },

  emits: ['refresh', 'create'],
};
</script>
