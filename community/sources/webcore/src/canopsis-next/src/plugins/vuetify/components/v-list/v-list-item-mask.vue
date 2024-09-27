<template>
  <span>{{ texts.start }}<span class="v-list-item__mask">{{ texts.middle }}</span>{{ texts.end }}</span>
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
      if (!props.mask) {
        return {
          start: props.text,
          middle: '',
          end: '',
        };
      }

      const index = props.text.toLocaleLowerCase().indexOf(props.mask.toLowerCase());

      if (index < 0) {
        return {
          start: props.text,
          middle: '',
          end: '',
        };
      }

      const middleEndIndex = index + props.mask.length;

      return {
        start: props.text.slice(0, index),
        middle: props.text.slice(index, middleEndIndex),
        end: props.text.slice(middleEndIndex),
      };
    });

    return {
      texts,
    };
  },
};
</script>
