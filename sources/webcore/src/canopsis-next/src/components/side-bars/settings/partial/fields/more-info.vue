<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.moreInfosModal')}}
    text-editor(:value="value", @input="$emit('input', $event)")
</template>

<script>
import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  components: {
    TextEditor,
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
};
</script>

