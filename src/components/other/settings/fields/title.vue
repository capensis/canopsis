<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('common.title')}}
    v-container
      v-text-field(:value="value", @input="$emit('input', $event)", :placeholder="$t('settings.widgetTitle')")
</template>

<script>
export default {
  props: {
    value: String,
  },
};
</script>
