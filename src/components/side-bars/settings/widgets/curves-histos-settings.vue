<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title
    v-btn(@click="", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import FieldTitle from '../partial/fields/title.vue';

export default {
  components: {
    FieldTitle,
  },
};
</script>

