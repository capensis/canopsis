<template>
  <c-alarm-chip v-bind="$attrs" v-on="$listeners">
    <template v-if="appendIconName" #append>
      <v-icon color="white" size="14">
        {{ appendIconName }}
      </v-icon>
    </template>
  </c-alarm-chip>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    appendIconName: {
      type: String,
      required: false,
    },
  },
};
</script>
