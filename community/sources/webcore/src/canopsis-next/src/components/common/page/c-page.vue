<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <slot />
    </v-card>
    <c-fab-btn
      :has-access="creatable"
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
    creatable: {
      type: Boolean,
      default: false,
    },
    createTooltip: {
      type: String,
      required: false,
    },
  },
};
</script>
