<template>
  <span>
    {{ texts.start }}<span class="v-list-item__mask">{{ texts.middle }}</span>{{ texts.end }}
  </span>
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    mask: {
      type: String,
      default: '',
    },
    text: {
      type: String,
      default: '',
    },
  },

  setup(props) {
    const texts = computed(() => {
      const text = String(props.text ?? '');
      const mask = String(props.mask ?? '');

      if (!mask) {
        return {
          start: text,
          middle: '',
          end: '',
        };
      }

      const index = text.toLocaleLowerCase().indexOf(mask.toLowerCase());

      if (index < 0) {
        return {
          start: text,
          middle: '',
          end: '',
        };
      }

      const middleEndIndex = index + mask.length;

      return {
        start: text.slice(0, index),
        middle: text.slice(index, middleEndIndex),
        end: text.slice(middleEndIndex),
      };
    });

    return {
      texts,
    };
  },
};
</script>

<style lang="scss">
.v-list-item__mask {
  color: rgb(var(--v-theme-primary));
}
</style>
