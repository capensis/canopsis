<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.yesNoMode')}}
    v-container
      v-switch(
      :label="$t('settings.yesNoMode')",
      :input-value="value",
      @change="$emit('input', $event)",
      color="primary"
      )
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
};
</script>
