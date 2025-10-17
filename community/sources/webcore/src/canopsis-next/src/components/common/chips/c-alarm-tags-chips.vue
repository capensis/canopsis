<template>
  <c-alarm-actions-chips
    :items="filteredTags"
    :active-items="selectedTags"
    :small="small"
    :inline-count="inlineCount"
    :closable-active="closableActive"
    item-class="c-alarm-tags-chips__chip"
    item-text="value"
    item-value="value"
    row
    v-on="$listeners"
  />
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    selectedTags: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
    small: {
      type: Boolean,
      default: false,
    },
    closableActive: {
      type: Boolean,
      default: false,
    },
    nameFilter: {
      type: String,
      default: '',
    },
    regexFilter: {
      type: String,
      default: '',
    },
    regexFilterFlags: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const filteredTags = computed(() => {
      const tags = (props.alarm.tag_colors ?? []);
      const regexps = [];

      if (props.nameFilter) {
        regexps.push(new RegExp(props.nameFilter));
      }

      if (props.regexFilter) {
        regexps.push(new RegExp(props.regexFilter, props.regexFilterFlags));
      }

      if (regexps.length) {
        return tags.filter(tag => regexps.every(regex => tag?.value?.match?.(regex)));
      }

      return tags;
    });

    return {
      filteredTags,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-tags-chips__chip .v-chip__content {
  padding: 0 4px;

  * {
    color: white;
  }
}
</style>
