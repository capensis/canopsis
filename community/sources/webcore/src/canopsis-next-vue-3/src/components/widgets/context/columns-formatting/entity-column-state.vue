<template>
  <c-alarm-chip v-bind="$attrs" @click="$emit('click', $event)">
    <template v-if="appendIconName" #append>
      <v-icon class="color--inherit" size="14">
        {{ appendIconName }}
      </v-icon>
    </template>
  </c-alarm-chip>
</template>

<script>
export default {
  inheritAttrs: false,

  emits: ['click'],

  props: {
    appendIconName: {
      type: String,
      required: false,
    },
  },
};
</script>
