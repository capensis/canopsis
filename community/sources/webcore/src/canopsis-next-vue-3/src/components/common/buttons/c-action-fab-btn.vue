<script setup>
defineOptions({
  inheritAttrs: false,
});

defineProps({
  icon: {
    type: String,
    required: true,
  },
  color: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['click']);
</script>

<template>
  <div class="c-action-fab-btn">
    <c-action-btn v-bind="$attrs">
      <template #button="{ props: buttonProps }">
        <v-btn
          v-bind="buttonProps"
          :color="color"
          size="small"
          class="c-small-fab"
          icon
          @click.stop="emit('click', $event)"
        >
          <v-icon :icon="icon" size="small" />
        </v-btn>
      </template>
    </c-action-btn>
  </div>
</template>

<style lang="scss" scoped>
.c-action-fab-btn {
  display: inline-flex;
}

.c-small-fab {
  --v-btn-height: 40px;

  width: 40px;
  height: 40px;
  border-radius: 50%;
}
</style>
